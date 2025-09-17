package vc

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Verifier handles the verification of Verifiable Credentials and Presentations
type Verifier struct {
	logger          *zap.Logger
	publicKeyStore  PublicKeyStore
	statusChecker   StatusChecker
	revocationStore RevocationStore
}

// PublicKeyStore interface for retrieving public keys
type PublicKeyStore interface {
	GetPublicKey(keyID string) (crypto.PublicKey, error)
}

// StatusChecker interface for checking credential status
type StatusChecker interface {
	CheckStatus(statusID string) (*CredentialStatusInfo, error)
}

// RevocationStore interface for checking revocation status
type RevocationStore interface {
	IsRevoked(credentialID string) (bool, error)
}

// DIDDocument represents a simplified DID document
type DIDDocument struct {
	ID                 string              `json:"id"`
	VerificationMethod []VerificationMethod `json:"verificationMethod"`
}

// VerificationMethod represents a verification method in a DID document
type VerificationMethod struct {
	ID           string      `json:"id"`
	Type         string      `json:"type"`
	Controller   string      `json:"controller"`
	PublicKeyPem string      `json:"publicKeyPem,omitempty"`
	PublicKeyJwk interface{} `json:"publicKeyJwk,omitempty"`
}

// CredentialStatusInfo represents credential status information
type CredentialStatusInfo struct {
	Valid   bool   `json:"valid"`
	Revoked bool   `json:"revoked"`
	Reason  string `json:"reason,omitempty"`
}

// VerificationResult represents the result of credential verification
type VerificationResult struct {
	Valid             bool                   `json:"valid"`
	Errors            []string               `json:"errors,omitempty"`
	Warnings          []string               `json:"warnings,omitempty"`
	CredentialStatus  *CredentialStatusInfo  `json:"credentialStatus,omitempty"`
	ProofVerification *ProofVerificationInfo `json:"proofVerification,omitempty"`
	Checks            map[string]bool        `json:"checks"`
}

// ProofVerificationInfo represents proof verification details
type ProofVerificationInfo struct {
	Valid             bool      `json:"valid"`
	ProofType         string    `json:"proofType"`
	VerificationMethod string    `json:"verificationMethod"`
	Created           time.Time `json:"created"`
	Purpose           string    `json:"purpose"`
}

// NewVerifier creates a new VC verifier
func NewVerifier(logger *zap.Logger, publicKeyStore PublicKeyStore, statusChecker StatusChecker, revocationStore RevocationStore) *Verifier {
	return &Verifier{
		logger:          logger,
		publicKeyStore:  publicKeyStore,
		statusChecker:   statusChecker,
		revocationStore: revocationStore,
	}
}

// VerifyCredential verifies a Verifiable Credential
func (v *Verifier) VerifyCredential(credential *VerifiableCredential, options *VerificationOptions) (*VerificationResult, error) {
	v.logger.Info("Verifying credential", zap.String("id", credential.ID))

	result := &VerificationResult{
		Valid:  true,
		Checks: make(map[string]bool),
	}

	// Basic validation
	if err := ValidateCredential(credential); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Basic validation failed: %v", err))
		result.Checks["basic_validation"] = false
	} else {
		result.Checks["basic_validation"] = true
	}

	// Verify proof if present
	if credential.Proof != nil {
		proofResult, err := v.verifyProof(credential.Proof, credential, options)
		if err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Proof verification failed: %v", err))
			result.Checks["proof_verification"] = false
		} else {
			result.ProofVerification = proofResult
			result.Checks["proof_verification"] = proofResult.Valid
			if !proofResult.Valid {
				result.Valid = false
			}
		}
	} else {
		result.Warnings = append(result.Warnings, "No proof present in credential")
		result.Checks["proof_verification"] = false
	}

	// Check expiration
	if credential.ExpirationDate != nil && time.Now().After(*credential.ExpirationDate) {
		result.Valid = false
		result.Errors = append(result.Errors, "Credential has expired")
		result.Checks["expiration"] = false
	} else {
		result.Checks["expiration"] = true
	}

	// Check credential status
	if credential.CredentialStatus != nil && v.statusChecker != nil {
		statusInfo, err := v.statusChecker.CheckStatus(credential.CredentialStatus.ID)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Could not check credential status: %v", err))
			result.Checks["status_check"] = false
		} else {
			result.CredentialStatus = statusInfo
			result.Checks["status_check"] = statusInfo.Valid
			if !statusInfo.Valid {
				result.Valid = false
				result.Errors = append(result.Errors, fmt.Sprintf("Credential status invalid: %s", statusInfo.Reason))
			}
		}
	} else {
		result.Checks["status_check"] = true // No status to check
	}

	// Check revocation
	if v.revocationStore != nil {
		revoked, err := v.revocationStore.IsRevoked(credential.ID)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Could not check revocation status: %v", err))
			result.Checks["revocation_check"] = false
		} else {
			result.Checks["revocation_check"] = !revoked
			if revoked {
				result.Valid = false
				result.Errors = append(result.Errors, "Credential has been revoked")
			}
		}
	} else {
		result.Checks["revocation_check"] = true // No revocation store to check
	}

	v.logger.Info("Credential verification completed", 
		zap.String("id", credential.ID),
		zap.Bool("valid", result.Valid),
		zap.Int("errors", len(result.Errors)),
		zap.Int("warnings", len(result.Warnings)))

	return result, nil
}

// VerifyPresentation verifies a Verifiable Presentation
func (v *Verifier) VerifyPresentation(presentation *VerifiablePresentation, options *VerificationOptions) (*VerificationResult, error) {
	v.logger.Info("Verifying presentation", zap.String("id", presentation.ID))

	result := &VerificationResult{
		Valid:  true,
		Checks: make(map[string]bool),
	}

	// Basic validation
	if err := ValidatePresentation(presentation); err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Basic validation failed: %v", err))
		result.Checks["basic_validation"] = false
	} else {
		result.Checks["basic_validation"] = true
	}

	// Verify presentation proof if present
	if presentation.Proof != nil {
		proofResult, err := v.verifyProof(presentation.Proof, presentation, options)
		if err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("Presentation proof verification failed: %v", err))
			result.Checks["presentation_proof"] = false
		} else {
			result.ProofVerification = proofResult
			result.Checks["presentation_proof"] = proofResult.Valid
			if !proofResult.Valid {
				result.Valid = false
			}
		}
	} else {
		result.Checks["presentation_proof"] = true // Presentation proof is optional
	}

	// Verify each credential in the presentation
	credentialErrors := 0
	for i, credential := range presentation.VerifiableCredential {
		credResult, err := v.VerifyCredential(&credential, options)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Error verifying credential %d: %v", i, err))
			credentialErrors++
		} else if !credResult.Valid {
			result.Errors = append(result.Errors, fmt.Sprintf("Credential %d is invalid: %v", i, strings.Join(credResult.Errors, ", ")))
			credentialErrors++
		}
	}

	result.Checks["credentials_valid"] = credentialErrors == 0
	if credentialErrors > 0 {
		result.Valid = false
	}

	// Verify challenge if provided
	if options != nil && options.Challenge != "" {
		if presentation.Proof == nil || presentation.Proof.Challenge != options.Challenge {
			result.Valid = false
			result.Errors = append(result.Errors, "Challenge mismatch")
			result.Checks["challenge"] = false
		} else {
			result.Checks["challenge"] = true
		}
	} else {
		result.Checks["challenge"] = true
	}

	// Verify domain if provided
	if options != nil && options.Domain != "" {
		if presentation.Proof == nil || presentation.Proof.Domain != options.Domain {
			result.Valid = false
			result.Errors = append(result.Errors, "Domain mismatch")
			result.Checks["domain"] = false
		} else {
			result.Checks["domain"] = true
		}
	} else {
		result.Checks["domain"] = true
	}

	v.logger.Info("Presentation verification completed", 
		zap.String("id", presentation.ID),
		zap.Bool("valid", result.Valid),
		zap.Int("errors", len(result.Errors)),
		zap.Int("warnings", len(result.Warnings)))

	return result, nil
}

// verifyProof verifies a cryptographic proof
func (v *Verifier) verifyProof(proof *Proof, document interface{}, options *VerificationOptions) (*ProofVerificationInfo, error) {
	proofInfo := &ProofVerificationInfo{
		ProofType:          proof.Type,
		VerificationMethod: proof.VerificationMethod,
		Created:            proof.Created,
		Purpose:            proof.ProofPurpose,
	}

	// Verify proof type
	if proof.Type != JSONWebSignature2020 {
		return proofInfo, fmt.Errorf("unsupported proof type: %s", proof.Type)
	}

	// Get public key for verification
	publicKey, err := v.getPublicKeyForVerification(proof.VerificationMethod)
	if err != nil {
		return proofInfo, fmt.Errorf("failed to get public key: %w", err)
	}

	// Create document copy without proof for verification
	documentCopy := v.removeProofFromDocument(document)
	documentBytes, err := json.Marshal(documentCopy)
	if err != nil {
		return proofInfo, fmt.Errorf("failed to marshal document: %w", err)
	}

	// Verify JWS
	if proof.JWS == "" {
		return proofInfo, fmt.Errorf("JWS is required for JsonWebSignature2020")
	}

	valid, err := v.verifyJWS(proof.JWS, documentBytes, publicKey)
	if err != nil {
		return proofInfo, fmt.Errorf("JWS verification failed: %w", err)
	}

	proofInfo.Valid = valid
	return proofInfo, nil
}

// getPublicKey retrieves the RSA public key for verification
func (v *Verifier) getPublicKey(verificationMethod string) (*rsa.PublicKey, error) {
	// Try to get directly from public key store
	pubKey, err := v.publicKeyStore.GetPublicKey(verificationMethod)
	if err == nil {
		rsaPubKey, ok := pubKey.(*rsa.PublicKey)
		if ok {
			return rsaPubKey, nil
		}
	}

	return nil, fmt.Errorf("public key not found for verification method: %s", verificationMethod)
}

// getPublicKeyForVerification retrieves the public key for verification
func (v *Verifier) getPublicKeyForVerification(verificationMethod string) (crypto.PublicKey, error) {
	if v.publicKeyStore == nil {
		return nil, fmt.Errorf("no public key store configured")
	}

	// Extract key ID from verification method
	parts := strings.Split(verificationMethod, "#")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid verification method format: %s", verificationMethod)
	}

	keyID := parts[1]

	// Try to get public key directly
	publicKey, err := v.publicKeyStore.GetPublicKey(keyID)
	if err == nil {
		return publicKey, nil
	}

	return nil, fmt.Errorf("public key not found for verification method: %s", verificationMethod)
}

// removeProofFromDocument removes proof from document for verification
func (v *Verifier) removeProofFromDocument(document interface{}) interface{} {
	switch doc := document.(type) {
	case *VerifiableCredential:
		copy := *doc
		copy.Proof = nil
		return copy
	case *VerifiablePresentation:
		copy := *doc
		copy.Proof = nil
		return copy
	default:
		return document
	}
}

// verifyJWS verifies a JSON Web Signature
func (v *Verifier) verifyJWS(jws string, payload []byte, publicKey crypto.PublicKey) (bool, error) {
	// Parse JWS (format: header..signature)
	parts := strings.Split(jws, ".")
	if len(parts) != 3 || parts[1] != "" {
		return false, fmt.Errorf("invalid JWS format")
	}

	headerB64 := parts[0]
	signatureB64 := parts[2]

	// Decode signature
	signature, err := base64.RawURLEncoding.DecodeString(signatureB64)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %w", err)
	}

	// Create signing input
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)
	signingInput := fmt.Sprintf("%s.%s", headerB64, payloadB64)

	// Verify signature
	return v.verifySignature([]byte(signingInput), signature, publicKey)
}

// verifySignature verifies a signature using the public key
func (v *Verifier) verifySignature(data, signature []byte, publicKey crypto.PublicKey) (bool, error) {
	switch key := publicKey.(type) {
	case *rsa.PublicKey:
		hash := sha256.Sum256(data)
		err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hash[:], signature)
		return err == nil, nil
	default:
		return false, fmt.Errorf("unsupported public key type: %T", publicKey)
	}
}

// parsePublicKeyFromPEM parses a public key from PEM format
func parsePublicKeyFromPEM(pemData []byte) (crypto.PublicKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	switch block.Type {
	case "PUBLIC KEY":
		return x509.ParsePKIXPublicKey(block.Bytes)
	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported public key type: %s", block.Type)
	}
}