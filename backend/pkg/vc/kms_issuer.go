package vc

import (
	"context"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/verza/pkg/kms"
)

// IssueCredentialRequest represents a request to issue a credential
type IssueCredentialRequest struct {
	SubjectDID        string      `json:"subjectDid"`
	CredentialTypes   []string    `json:"credentialTypes"`
	CredentialSubject interface{} `json:"credentialSubject"`
	ExpirationDate    *time.Time  `json:"expirationDate,omitempty"`
	CredentialID      string      `json:"credentialId,omitempty"`
	Context           []string    `json:"@context,omitempty"`
	Challenge         string      `json:"challenge,omitempty"`
	Domain            string      `json:"domain,omitempty"`
	StatusListURL     string      `json:"statusListUrl,omitempty"`
	StatusListIndex   string      `json:"statusListIndex,omitempty"`
}

// KMSIssuer handles the creation and signing of Verifiable Credentials using KMS
type KMSIssuer struct {
	logger    *zap.Logger
	issuerDID string
	kms       kms.KMS
	keyID     string
	algorithm kms.SigningAlgorithm
}

// NewKMSIssuer creates a new VC issuer that uses KMS for signing
func NewKMSIssuer(logger *zap.Logger, issuerDID string, kmsInstance kms.KMS, keyID string, algorithm kms.SigningAlgorithm) (*KMSIssuer, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if issuerDID == "" {
		return nil, fmt.Errorf("issuer DID is required")
	}
	if kmsInstance == nil {
		return nil, fmt.Errorf("KMS instance is required")
	}
	if keyID == "" {
		return nil, fmt.Errorf("key ID is required")
	}

	// Verify key exists in KMS
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := kmsInstance.GetKeyInfo(ctx, keyID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify key in KMS: %w", err)
	}

	return &KMSIssuer{
		logger:    logger,
		issuerDID: issuerDID,
		kms:       kmsInstance,
		keyID:     keyID,
		algorithm: algorithm,
	}, nil
}

// IssueCredential creates and signs a new Verifiable Credential using KMS
func (k *KMSIssuer) IssueCredential(ctx context.Context, req IssueCredentialRequest) (*VerifiableCredential, error) {
	k.logger.Info("Issuing credential with KMS", 
		zap.String("subject_did", req.SubjectDID),
		zap.String("key_id", k.keyID),
		zap.Strings("types", req.CredentialTypes))

	// Validate request
	if err := k.validateIssueRequest(req); err != nil {
		return nil, fmt.Errorf("invalid issue request: %w", err)
	}

	// Create credential ID
	credentialID := fmt.Sprintf("urn:uuid:%s", uuid.New().String())

	// Set default expiration if not provided
	var expirationDate *time.Time
	if req.ExpirationDate != nil {
		expirationDate = req.ExpirationDate
	} else {
		defaultExp := time.Now().AddDate(1, 0, 0)
		expirationDate = &defaultExp
	}

	// Create the credential
	credential := &VerifiableCredential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://verza.io/contexts/v1",
		},
		ID:               credentialID,
		Type:             append([]string{"VerifiableCredential"}, req.CredentialTypes...),
		Issuer:           CredentialIssuer{ID: k.issuerDID},
		IssuanceDate:     time.Now(),
		ExpirationDate:   expirationDate,
		CredentialSubject: req.CredentialSubject,
	}

	// Add credential status if provided
	if req.StatusListURL != "" {
		credential.CredentialStatus = &CredentialStatus{
			ID:   req.StatusListURL,
			Type: "StatusList2021Entry",
		}
	}

	// Create proof using KMS
	proof, err := k.createProofWithKMS(ctx, credential, req.Challenge, req.Domain)
	if err != nil {
		return nil, fmt.Errorf("failed to create proof: %w", err)
	}

	credential.Proof = proof

	k.logger.Info("Credential issued successfully", 
		zap.String("credential_id", credentialID),
		zap.String("subject_did", req.SubjectDID))

	return credential, nil
}

// GetPublicKey returns the public key for this issuer from KMS
func (k *KMSIssuer) GetPublicKey(ctx context.Context) (crypto.PublicKey, error) {
	return k.kms.GetPublicKey(ctx, k.keyID)
}

// GetKeyID returns the key ID used by this issuer
func (k *KMSIssuer) GetKeyID() string {
	return k.keyID
}

// GetIssuerDID returns the issuer DID
func (k *KMSIssuer) GetIssuerDID() string {
	return k.issuerDID
}

// RotateKey rotates the signing key in KMS
func (k *KMSIssuer) RotateKey(ctx context.Context) error {
	k.logger.Info("Rotating issuer key", zap.String("key_id", k.keyID))

	err := k.kms.RotateKey(ctx, k.keyID)
	if err != nil {
		return fmt.Errorf("failed to rotate key: %w", err)
	}

	k.logger.Info("Key rotated successfully", zap.String("key_id", k.keyID))
	return nil
}

// Helper methods

func (k *KMSIssuer) validateIssueRequest(req IssueCredentialRequest) error {
	if req.SubjectDID == "" {
		return fmt.Errorf("subject DID is required")
	}
	if len(req.CredentialTypes) == 0 {
		return fmt.Errorf("at least one credential type is required")
	}
	if req.CredentialSubject == nil {
		return fmt.Errorf("credential subject is required")
	}
	return nil
}

func (k *KMSIssuer) createProofWithKMS(ctx context.Context, credential *VerifiableCredential, challenge, domain string) (*Proof, error) {
	// Create proof metadata
	proof := &Proof{
		Type:               "JsonWebSignature2020",
		Created:            time.Now(),
		VerificationMethod: fmt.Sprintf("%s#%s", k.issuerDID, k.keyID),
		ProofPurpose:       "assertionMethod",
	}

	// Add challenge and domain if provided
	if challenge != "" {
		proof.Challenge = challenge
	}
	if domain != "" {
		proof.Domain = domain
	}

	// Create JWS using KMS
	jws, err := k.createJWSWithKMS(ctx, credential, proof)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWS: %w", err)
	}

	proof.JWS = jws
	return proof, nil
}

func (k *KMSIssuer) createJWSWithKMS(ctx context.Context, credential *VerifiableCredential, proof *Proof) (string, error) {
	// Create JWT header
	header := map[string]interface{}{
		"alg": string(k.algorithm),
		"typ": "JWT",
		"kid": k.keyID,
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}

	// Create JWT payload (credential without proof)
	credentialCopy := *credential
	credentialCopy.Proof = nil // Remove proof for signing

	payloadBytes, err := json.Marshal(credentialCopy)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Encode header and payload
	headerB64 := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// Create signing input
	signingInput := fmt.Sprintf("%s.%s", headerB64, payloadB64)

	// Sign with KMS
	signReq := kms.SignRequest{
		KeyID:     k.keyID,
		Data:      []byte(signingInput),
		Algorithm: k.algorithm,
		Context: map[string]string{
			"purpose": "credential_signing",
			"issuer":  k.issuerDID,
		},
	}

	signResp, err := k.kms.Sign(ctx, signReq)
	if err != nil {
		return "", fmt.Errorf("failed to sign with KMS: %w", err)
	}

	// Encode signature
	signatureB64 := base64.RawURLEncoding.EncodeToString(signResp.Signature)

	// Return complete JWS
	return fmt.Sprintf("%s.%s", signingInput, signatureB64), nil
}

// CreateJWT creates a JWT token using KMS for external authentication
func (k *KMSIssuer) CreateJWT(ctx context.Context, claims jwt.MapClaims, audience string) (string, error) {
	k.logger.Debug("Creating JWT with KMS", zap.String("audience", audience))

	// Set standard claims
	now := time.Now()
	claims["iss"] = k.issuerDID
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(time.Hour).Unix() // 1 hour expiration
	if audience != "" {
		claims["aud"] = audience
	}
	if _, exists := claims["jti"]; !exists {
		claims["jti"] = uuid.New().String()
	}

	// Create JWT header
	header := map[string]interface{}{
		"alg": string(k.algorithm),
		"typ": "JWT",
		"kid": k.keyID,
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Encode header and payload
	headerB64 := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// Create signing input
	signingInput := fmt.Sprintf("%s.%s", headerB64, payloadB64)

	// Sign with KMS
	signReq := kms.SignRequest{
		KeyID:     k.keyID,
		Data:      []byte(signingInput),
		Algorithm: k.algorithm,
		Context: map[string]string{
			"purpose": "jwt_signing",
			"issuer":  k.issuerDID,
			"audience": audience,
		},
	}

	signResp, err := k.kms.Sign(ctx, signReq)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT with KMS: %w", err)
	}

	// Encode signature
	signatureB64 := base64.RawURLEncoding.EncodeToString(signResp.Signature)

	// Return complete JWT
	return fmt.Sprintf("%s.%s", signingInput, signatureB64), nil
}