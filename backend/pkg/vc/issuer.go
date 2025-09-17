package vc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Issuer handles the creation and signing of Verifiable Credentials
type Issuer struct {
	logger     *zap.Logger
	issuerDID  string
	privateKey crypto.PrivateKey
	publicKey  crypto.PublicKey
	keyID      string
}

// NewIssuer creates a new VC issuer
func NewIssuer(logger *zap.Logger, issuerDID string, privateKey crypto.PrivateKey, keyID string) (*Issuer, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if issuerDID == "" {
		return nil, fmt.Errorf("issuer DID is required")
	}
	if privateKey == nil {
		return nil, fmt.Errorf("private key is required")
	}
	if keyID == "" {
		return nil, fmt.Errorf("key ID is required")
	}

	var publicKey crypto.PublicKey
	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		publicKey = &key.PublicKey
	default:
		return nil, fmt.Errorf("unsupported private key type: %T", privateKey)
	}

	return &Issuer{
		logger:     logger,
		issuerDID:  issuerDID,
		privateKey: privateKey,
		publicKey:  publicKey,
		keyID:      keyID,
	}, nil
}

// IssueCredential creates and signs a new Verifiable Credential
func (i *Issuer) IssueCredential(subject CredentialSubject, credentialType string, options *SigningOptions) (*VerifiableCredential, error) {
	i.logger.Info("Issuing new credential", zap.String("type", credentialType), zap.String("subject", subject.ID))

	// Set default options if not provided
	if options == nil {
		options = &SigningOptions{
			VerificationMethod: fmt.Sprintf("%s#%s", i.issuerDID, i.keyID),
			ProofPurpose:       AssertionMethod,
			Created:            time.Now().UTC(),
		}
	}

	// Create the credential
	credential := &VerifiableCredential{
		Context:           append(DefaultContexts(), "https://verza.io/contexts/v1"),
		ID:                fmt.Sprintf("urn:uuid:%s", uuid.New().String()),
		Type:              append(DefaultVCTypes(), credentialType),
		Issuer:            CredentialIssuer{ID: i.issuerDID},
		IssuanceDate:      time.Now().UTC(),
		CredentialSubject: subject,
	}

	// Add credential status if needed
	if credentialType != "TemporaryCredential" {
		credential.CredentialStatus = &CredentialStatus{
			ID:   fmt.Sprintf("https://verza.io/status/%s", credential.ID),
			Type: StatusList2021Entry,
		}
	}

	// Sign the credential
	if err := i.signCredential(credential, options); err != nil {
		i.logger.Error("Failed to sign credential", zap.Error(err))
		return nil, fmt.Errorf("failed to sign credential: %w", err)
	}

	i.logger.Info("Credential issued successfully", zap.String("id", credential.ID))
	return credential, nil
}

// signCredential adds a cryptographic proof to the credential
func (i *Issuer) signCredential(credential *VerifiableCredential, options *SigningOptions) error {
	// Create a copy of the credential without proof for signing
	credentialCopy := *credential
	credentialCopy.Proof = nil

	// Serialize the credential for signing
	credentialBytes, err := json.Marshal(credentialCopy)
	if err != nil {
		return fmt.Errorf("failed to marshal credential: %w", err)
	}

	// Create the proof
	proof := &Proof{
		Type:               JSONWebSignature2020,
		Created:            options.Created,
		VerificationMethod: options.VerificationMethod,
		ProofPurpose:       options.ProofPurpose,
		Challenge:          options.Challenge,
		Domain:             options.Domain,
	}

	// Create JWS
	jws, err := i.createJWS(credentialBytes, proof)
	if err != nil {
		return fmt.Errorf("failed to create JWS: %w", err)
	}

	proof.JWS = jws
	credential.Proof = proof

	return nil
}

// createJWS creates a JSON Web Signature for the credential
func (i *Issuer) createJWS(payload []byte, proof *Proof) (string, error) {
	// Create JWS header
	header := map[string]interface{}{
		"alg": "RS256",
		"typ": "JWT",
		"kid": i.keyID,
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}

	// Base64URL encode header and payload
	headerB64 := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)

	// Create signing input
	signingInput := fmt.Sprintf("%s.%s", headerB64, payloadB64)

	// Sign the input
	signature, err := i.signData([]byte(signingInput))
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}

	// Base64URL encode signature
	signatureB64 := base64.RawURLEncoding.EncodeToString(signature)

	// Return compact JWS
	return fmt.Sprintf("%s..%s", headerB64, signatureB64), nil
}

// signData signs the given data using the issuer's private key
func (i *Issuer) signData(data []byte) ([]byte, error) {
	switch key := i.privateKey.(type) {
	case *rsa.PrivateKey:
		hash := sha256.Sum256(data)
		return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	default:
		return nil, fmt.Errorf("unsupported private key type: %T", i.privateKey)
	}
}

// CreatePresentation creates a Verifiable Presentation from credentials
func (i *Issuer) CreatePresentation(credentials []VerifiableCredential, holder string, options *SigningOptions) (*VerifiablePresentation, error) {
	i.logger.Info("Creating presentation", zap.String("holder", holder), zap.Int("credentials", len(credentials)))

	// Set default options if not provided
	if options == nil {
		options = &SigningOptions{
			VerificationMethod: fmt.Sprintf("%s#%s", holder, "key-1"),
			ProofPurpose:       Authentication,
			Created:            time.Now().UTC(),
		}
	}

	// Create the presentation
	presentation := &VerifiablePresentation{
		Context:              DefaultContexts(),
		ID:                   fmt.Sprintf("urn:uuid:%s", uuid.New().String()),
		Type:                 DefaultVPTypes(),
		Holder:               holder,
		VerifiableCredential: credentials,
	}

	// Sign the presentation if holder matches issuer
	if holder == i.issuerDID {
		if err := i.signPresentation(presentation, options); err != nil {
			i.logger.Error("Failed to sign presentation", zap.Error(err))
			return nil, fmt.Errorf("failed to sign presentation: %w", err)
		}
	}

	i.logger.Info("Presentation created successfully", zap.String("id", presentation.ID))
	return presentation, nil
}

// signPresentation adds a cryptographic proof to the presentation
func (i *Issuer) signPresentation(presentation *VerifiablePresentation, options *SigningOptions) error {
	// Create a copy of the presentation without proof for signing
	presentationCopy := *presentation
	presentationCopy.Proof = nil

	// Serialize the presentation for signing
	presentationBytes, err := json.Marshal(presentationCopy)
	if err != nil {
		return fmt.Errorf("failed to marshal presentation: %w", err)
	}

	// Create the proof
	proof := &Proof{
		Type:               JSONWebSignature2020,
		Created:            options.Created,
		VerificationMethod: options.VerificationMethod,
		ProofPurpose:       options.ProofPurpose,
		Challenge:          options.Challenge,
		Domain:             options.Domain,
	}

	// Create JWS
	jws, err := i.createJWS(presentationBytes, proof)
	if err != nil {
		return fmt.Errorf("failed to create JWS: %w", err)
	}

	proof.JWS = jws
	presentation.Proof = proof

	return nil
}

// ParsePrivateKeyFromPEM parses a private key from PEM format
func ParsePrivateKeyFromPEM(pemData []byte) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		return x509.ParsePKCS8PrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported private key type: %s", block.Type)
	}
}

// GetIssuerDID returns the issuer's DID
func (i *Issuer) GetIssuerDID() string {
	return i.issuerDID
}

// GetKeyID returns the key ID
func (i *Issuer) GetKeyID() string {
	return i.keyID
}

// ValidateCredential performs basic validation on a credential
func ValidateCredential(credential *VerifiableCredential) error {
	if credential == nil {
		return fmt.Errorf("credential is nil")
	}

	if len(credential.Context) == 0 {
		return fmt.Errorf("credential context is required")
	}

	if len(credential.Type) == 0 {
		return fmt.Errorf("credential type is required")
	}

	if credential.Issuer.ID == "" {
		return fmt.Errorf("credential issuer is required")
	}

	if credential.IssuanceDate.IsZero() {
		return fmt.Errorf("credential issuance date is required")
	}

	if credential.CredentialSubject == nil {
		return fmt.Errorf("credential subject is required")
	}

	// Check if credential has expired
	if credential.ExpirationDate != nil && time.Now().After(*credential.ExpirationDate) {
		return fmt.Errorf("credential has expired")
	}

	// Validate that VerifiableCredential type is present
	found := false
	for _, t := range credential.Type {
		if t == VerifiableCredentialType {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("credential must include VerifiableCredential type")
	}

	return nil
}

// ValidatePresentation performs basic validation on a presentation
func ValidatePresentation(presentation *VerifiablePresentation) error {
	if presentation == nil {
		return fmt.Errorf("presentation is nil")
	}

	if len(presentation.Context) == 0 {
		return fmt.Errorf("presentation context is required")
	}

	if len(presentation.Type) == 0 {
		return fmt.Errorf("presentation type is required")
	}

	// Validate that VerifiablePresentation type is present
	found := false
	for _, t := range presentation.Type {
		if t == VerifiablePresentationType {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("presentation must include VerifiablePresentation type")
	}

	// Validate each credential in the presentation
	for i, credential := range presentation.VerifiableCredential {
		if err := ValidateCredential(&credential); err != nil {
			return fmt.Errorf("invalid credential at index %d: %w", i, err)
		}
	}

	return nil
}