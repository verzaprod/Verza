package kms

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// LocalKMS implements an in-memory KMS for development and testing
type LocalKMS struct {
	logger *zap.Logger
	keys   map[string]*localKey
	mutex  sync.RWMutex
}

type localKey struct {
	info       *KeyInfo
	privateKey crypto.PrivateKey
	publicKey  crypto.PublicKey
}

// NewLocalKMS creates a new local KMS instance
func NewLocalKMS(logger *zap.Logger) (*LocalKMS, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &LocalKMS{
		logger: logger,
		keys:   make(map[string]*localKey),
	}, nil
}

// CreateKey creates a new key in memory
func (l *LocalKMS) CreateKey(ctx context.Context, keyID string, keyType KeyType) (*KeyInfo, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.logger.Info("Creating local key", zap.String("key_id", keyID), zap.String("key_type", string(keyType)))

	// Check if key already exists
	if _, exists := l.keys[keyID]; exists {
		return nil, fmt.Errorf("key %s already exists", keyID)
	}

	// Generate key pair based on type
	privateKey, publicKey, algorithm, err := l.generateKeyPair(keyType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Create key info
	keyInfo := &KeyInfo{
		KeyID:     keyID,
		KeyType:   keyType,
		Algorithm: algorithm,
		PublicKey: publicKey,
		CreatedAt: time.Now(),
		Enabled:   true,
	}

	// Store key
	l.keys[keyID] = &localKey{
		info:       keyInfo,
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	return keyInfo, nil
}

// GetKeyInfo retrieves information about a key
func (l *LocalKMS) GetKeyInfo(ctx context.Context, keyID string) (*KeyInfo, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	key, exists := l.keys[keyID]
	if !exists {
		return nil, fmt.Errorf("key %s not found", keyID)
	}

	return key.info, nil
}

// GetPublicKey retrieves the public key for a given key ID
func (l *LocalKMS) GetPublicKey(ctx context.Context, keyID string) (crypto.PublicKey, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	key, exists := l.keys[keyID]
	if !exists {
		return nil, fmt.Errorf("key %s not found", keyID)
	}

	return key.publicKey, nil
}

// Sign signs data using the specified key
func (l *LocalKMS) Sign(ctx context.Context, req SignRequest) (*SignResponse, error) {
	l.mutex.RLock()
	key, exists := l.keys[req.KeyID]
	l.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("key %s not found", req.KeyID)
	}

	l.logger.Debug("Signing data locally", zap.String("key_id", req.KeyID), zap.String("algorithm", string(req.Algorithm)))

	// Sign based on key type
	signature, err := l.signData(key.privateKey, req.Data, req.Algorithm)
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}

	return &SignResponse{
		Signature: signature,
		KeyID:     req.KeyID,
	}, nil
}

// RotateKey rotates a key to a new version (recreates the key)
func (l *LocalKMS) RotateKey(ctx context.Context, keyID string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	key, exists := l.keys[keyID]
	if !exists {
		return fmt.Errorf("key %s not found", keyID)
	}

	l.logger.Info("Rotating local key", zap.String("key_id", keyID))

	// Generate new key pair
	privateKey, publicKey, algorithm, err := l.generateKeyPair(key.info.KeyType)
	if err != nil {
		return fmt.Errorf("failed to generate new key pair: %w", err)
	}

	// Update key
	key.privateKey = privateKey
	key.publicKey = publicKey
	key.info.PublicKey = publicKey
	key.info.Algorithm = algorithm
	key.info.CreatedAt = time.Now()

	return nil
}

// DeleteKey deletes a key from memory
func (l *LocalKMS) DeleteKey(ctx context.Context, keyID string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if _, exists := l.keys[keyID]; !exists {
		return fmt.Errorf("key %s not found", keyID)
	}

	l.logger.Info("Deleting local key", zap.String("key_id", keyID))
	delete(l.keys, keyID)

	return nil
}

// ListKeys lists all keys in memory
func (l *LocalKMS) ListKeys(ctx context.Context) ([]string, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	keys := make([]string, 0, len(l.keys))
	for keyID := range l.keys {
		keys = append(keys, keyID)
	}

	return keys, nil
}

// Helper methods

func (l *LocalKMS) generateKeyPair(keyType KeyType) (crypto.PrivateKey, crypto.PublicKey, SigningAlgorithm, error) {
	switch keyType {
	case KeyTypeRSA2048:
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, nil, "", err
		}
		return privateKey, &privateKey.PublicKey, AlgRS256, nil

	case KeyTypeRSA4096:
		privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, nil, "", err
		}
		return privateKey, &privateKey.PublicKey, AlgRS256, nil

	case KeyTypeECDSAP256:
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, nil, "", err
		}
		return privateKey, &privateKey.PublicKey, AlgES256, nil

	case KeyTypeECDSAP384:
		privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			return nil, nil, "", err
		}
		return privateKey, &privateKey.PublicKey, AlgES384, nil

	case KeyTypeECDSAP521:
		privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return nil, nil, "", err
		}
		return privateKey, &privateKey.PublicKey, AlgES512, nil

	case KeyTypeEd25519:
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil, "", err
		}
		return privateKey, publicKey, AlgEdDSA, nil

	default:
		return nil, nil, "", fmt.Errorf("unsupported key type: %s", keyType)
	}
}

func (l *LocalKMS) signData(privateKey crypto.PrivateKey, data []byte, algorithm SigningAlgorithm) ([]byte, error) {
	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		return l.signRSA(key, data, algorithm)
	case *ecdsa.PrivateKey:
		return l.signECDSA(key, data, algorithm)
	case ed25519.PrivateKey:
		return l.signEd25519(key, data)
	default:
		return nil, fmt.Errorf("unsupported private key type: %T", privateKey)
	}
}

func (l *LocalKMS) signRSA(privateKey *rsa.PrivateKey, data []byte, algorithm SigningAlgorithm) ([]byte, error) {
	hash := sha256.Sum256(data)

	switch algorithm {
	case AlgRS256, AlgRS384, AlgRS512:
		return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	default:
		return nil, fmt.Errorf("unsupported RSA algorithm: %s", algorithm)
	}
}

func (l *LocalKMS) signECDSA(privateKey *ecdsa.PrivateKey, data []byte, algorithm SigningAlgorithm) ([]byte, error) {
	hash := sha256.Sum256(data)

	switch algorithm {
	case AlgES256, AlgES384, AlgES512:
		return ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	default:
		return nil, fmt.Errorf("unsupported ECDSA algorithm: %s", algorithm)
	}
}

func (l *LocalKMS) signEd25519(privateKey ed25519.PrivateKey, data []byte) ([]byte, error) {
	return ed25519.Sign(privateKey, data), nil
}