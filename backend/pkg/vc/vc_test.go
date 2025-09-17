package vc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Mock implementations for testing
type mockPublicKeyStore struct {
	keys map[string]interface{}
	dids map[string]*DIDDocument
}

func (m *mockPublicKeyStore) GetPublicKey(keyID string) (crypto.PublicKey, error) {
	if key, exists := m.keys[keyID]; exists {
		return key, nil
	}
	return nil, fmt.Errorf("key not found: %s", keyID)
}

func (m *mockPublicKeyStore) GetDIDDocument(did string) (*DIDDocument, error) {
	if doc, exists := m.dids[did]; exists {
		return doc, nil
	}
	return nil, fmt.Errorf("DID document not found: %s", did)
}

type mockStatusChecker struct{}

func (m *mockStatusChecker) CheckStatus(statusID string) (*CredentialStatusInfo, error) {
	return &CredentialStatusInfo{
		Valid:   true,
		Revoked: false,
	}, nil
}

type mockRevocationStore struct {
	revokedCredentials map[string]bool
}

func (m *mockRevocationStore) IsRevoked(credentialID string) (bool, error) {
	return m.revokedCredentials[credentialID], nil
}

func TestVerifiableCredentialTypes(t *testing.T) {
	// Test CredentialSubject JSON marshaling/unmarshaling
	subject := CredentialSubject{
		ID: "did:example:123",
		Data: map[string]interface{}{
			"name":  "John Doe",
			"email": "john@example.com",
			"age":   30,
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(subject)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaled CredentialSubject
	err = json.Unmarshal(jsonData, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, subject.ID, unmarshaled.ID)
	assert.Equal(t, subject.Data["name"], unmarshaled.Data["name"])
	assert.Equal(t, subject.Data["email"], unmarshaled.Data["email"])
	assert.Equal(t, float64(30), unmarshaled.Data["age"]) // JSON numbers are float64
}

func TestIssuerCreation(t *testing.T) {
	logger := zap.NewNop()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	issuerDID := "did:example:issuer"
	keyID := "key-1"

	// Test successful issuer creation
	issuer, err := NewIssuer(logger, issuerDID, privateKey, keyID)
	require.NoError(t, err)
	assert.NotNil(t, issuer)
	assert.Equal(t, issuerDID, issuer.GetIssuerDID())
	assert.Equal(t, keyID, issuer.GetKeyID())

	// Test error cases
	_, err = NewIssuer(nil, issuerDID, privateKey, keyID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "logger is required")

	_, err = NewIssuer(logger, "", privateKey, keyID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "issuer DID is required")

	_, err = NewIssuer(logger, issuerDID, nil, keyID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "private key is required")

	_, err = NewIssuer(logger, issuerDID, privateKey, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key ID is required")
}

func TestCredentialIssuance(t *testing.T) {
	logger := zap.NewNop()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	issuerDID := "did:example:issuer"
	keyID := "key-1"

	issuer, err := NewIssuer(logger, issuerDID, privateKey, keyID)
	require.NoError(t, err)

	// Create credential subject
	subject := CredentialSubject{
		ID: "did:example:subject",
		Data: map[string]interface{}{
			"name":       "John Doe",
			"birthDate":  "1990-01-01",
			"nationality": "US",
		},
	}

	// Issue credential
	credential, err := issuer.IssueCredential(subject, "IdentityCredential", nil)
	require.NoError(t, err)
	assert.NotNil(t, credential)

	// Verify credential structure
	assert.NotEmpty(t, credential.ID)
	assert.Contains(t, credential.Context, W3CCredentialsContext)
	assert.Contains(t, credential.Type, VerifiableCredentialType)
	assert.Contains(t, credential.Type, "IdentityCredential")
	assert.Equal(t, issuerDID, credential.Issuer.ID)
	assert.NotNil(t, credential.Proof)
	assert.Equal(t, JSONWebSignature2020, credential.Proof.Type)
	assert.NotEmpty(t, credential.Proof.JWS)

	// Verify credential subject
	credSubject, ok := credential.CredentialSubject.(CredentialSubject)
	require.True(t, ok)
	assert.Equal(t, subject.ID, credSubject.ID)
	assert.Equal(t, subject.Data["name"], credSubject.Data["name"])
}

func TestPresentationCreation(t *testing.T) {
	logger := zap.NewNop()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	issuerDID := "did:example:issuer"
	keyID := "key-1"

	issuer, err := NewIssuer(logger, issuerDID, privateKey, keyID)
	require.NoError(t, err)

	// Create and issue a credential
	subject := CredentialSubject{
		ID: "did:example:subject",
		Data: map[string]interface{}{
			"name": "John Doe",
		},
	}

	credential, err := issuer.IssueCredential(subject, "IdentityCredential", nil)
	require.NoError(t, err)

	// Create presentation
	credentials := []VerifiableCredential{*credential}
	holder := issuerDID // Same as issuer for this test

	presentation, err := issuer.CreatePresentation(credentials, holder, nil)
	require.NoError(t, err)
	assert.NotNil(t, presentation)

	// Verify presentation structure
	assert.NotEmpty(t, presentation.ID)
	assert.Contains(t, presentation.Context, W3CCredentialsContext)
	assert.Contains(t, presentation.Type, VerifiablePresentationType)
	assert.Equal(t, holder, presentation.Holder)
	assert.Len(t, presentation.VerifiableCredential, 1)
	assert.NotNil(t, presentation.Proof) // Should have proof since holder == issuer
}

func TestCredentialValidation(t *testing.T) {
	// Test valid credential
	validCredential := &VerifiableCredential{
		Context:      DefaultContexts(),
		ID:           "urn:uuid:test",
		Type:         DefaultVCTypes(),
		Issuer:       CredentialIssuer{ID: "did:example:issuer"},
		IssuanceDate: time.Now(),
		CredentialSubject: CredentialSubject{
			ID: "did:example:subject",
			Data: map[string]interface{}{
				"name": "Test",
			},
		},
	}

	err := ValidateCredential(validCredential)
	assert.NoError(t, err)

	// Test nil credential
	err = ValidateCredential(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "credential is nil")

	// Test missing context
	invalidCredential := *validCredential
	invalidCredential.Context = nil
	err = ValidateCredential(&invalidCredential)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context is required")

	// Test missing type
	invalidCredential = *validCredential
	invalidCredential.Type = nil
	err = ValidateCredential(&invalidCredential)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type is required")

	// Test missing issuer
	invalidCredential = *validCredential
	invalidCredential.Issuer = CredentialIssuer{}
	err = ValidateCredential(&invalidCredential)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "issuer is required")

	// Test expired credential
	expiredCredential := *validCredential
	expiredDate := time.Now().Add(-24 * time.Hour)
	expiredCredential.ExpirationDate = &expiredDate
	err = ValidateCredential(&expiredCredential)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "has expired")

	// Test missing VerifiableCredential type
	invalidCredential = *validCredential
	invalidCredential.Type = []string{"CustomType"}
	err = ValidateCredential(&invalidCredential)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must include VerifiableCredential type")
}

func TestPresentationValidation(t *testing.T) {
	// Create a valid credential first
	validCredential := VerifiableCredential{
		Context:      DefaultContexts(),
		ID:           "urn:uuid:test",
		Type:         DefaultVCTypes(),
		Issuer:       CredentialIssuer{ID: "did:example:issuer"},
		IssuanceDate: time.Now(),
		CredentialSubject: CredentialSubject{
			ID: "did:example:subject",
			Data: map[string]interface{}{
				"name": "Test",
			},
		},
	}

	// Test valid presentation
	validPresentation := &VerifiablePresentation{
		Context:              DefaultContexts(),
		ID:                   "urn:uuid:presentation",
		Type:                 DefaultVPTypes(),
		Holder:               "did:example:holder",
		VerifiableCredential: []VerifiableCredential{validCredential},
	}

	err := ValidatePresentation(validPresentation)
	assert.NoError(t, err)

	// Test nil presentation
	err = ValidatePresentation(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "presentation is nil")

	// Test missing VerifiablePresentation type
	invalidPresentation := *validPresentation
	invalidPresentation.Type = []string{"CustomType"}
	err = ValidatePresentation(&invalidPresentation)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must include VerifiablePresentation type")
}

func TestVerifierCreation(t *testing.T) {
	logger := zap.NewNop()
	publicKeyStore := &mockPublicKeyStore{
		keys: make(map[string]interface{}),
		dids: make(map[string]*DIDDocument),
	}
	statusChecker := &mockStatusChecker{}
	revocationStore := &mockRevocationStore{
		revokedCredentials: make(map[string]bool),
	}

	verifier := NewVerifier(logger, publicKeyStore, statusChecker, revocationStore)
	assert.NotNil(t, verifier)
}

func TestCredentialVerification(t *testing.T) {
	logger := zap.NewNop()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	issuerDID := "did:example:issuer"
	keyID := "key-1"

	// Create issuer and issue credential
	issuer, err := NewIssuer(logger, issuerDID, privateKey, keyID)
	require.NoError(t, err)

	subject := CredentialSubject{
		ID: "did:example:subject",
		Data: map[string]interface{}{
			"name": "John Doe",
		},
	}

	credential, err := issuer.IssueCredential(subject, "IdentityCredential", nil)
	require.NoError(t, err)

	// Set up verifier with mock stores
	publicKeyStore := &mockPublicKeyStore{
		keys: map[string]interface{}{
			keyID: &privateKey.PublicKey,
		},
		dids: make(map[string]*DIDDocument),
	}
	statusChecker := &mockStatusChecker{}
	revocationStore := &mockRevocationStore{
		revokedCredentials: make(map[string]bool),
	}

	verifier := NewVerifier(logger, publicKeyStore, statusChecker, revocationStore)

	// Verify credential
	result, err := verifier.VerifyCredential(credential, nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid)
	assert.Empty(t, result.Errors)
	assert.True(t, result.Checks["basic_validation"])
	assert.True(t, result.Checks["expiration"])
	assert.True(t, result.Checks["status_check"])
	assert.True(t, result.Checks["revocation_check"])
}

func TestPresentationVerification(t *testing.T) {
	logger := zap.NewNop()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	issuerDID := "did:example:issuer"
	keyID := "key-1"

	// Create issuer and issue credential
	issuer, err := NewIssuer(logger, issuerDID, privateKey, keyID)
	require.NoError(t, err)

	subject := CredentialSubject{
		ID: "did:example:subject",
		Data: map[string]interface{}{
			"name": "John Doe",
		},
	}

	credential, err := issuer.IssueCredential(subject, "IdentityCredential", nil)
	require.NoError(t, err)

	// Create presentation
	credentials := []VerifiableCredential{*credential}
	presentation, err := issuer.CreatePresentation(credentials, issuerDID, nil)
	require.NoError(t, err)

	// Set up verifier
	publicKeyStore := &mockPublicKeyStore{
		keys: map[string]interface{}{
			keyID: &privateKey.PublicKey,
		},
		dids: make(map[string]*DIDDocument),
	}
	statusChecker := &mockStatusChecker{}
	revocationStore := &mockRevocationStore{
		revokedCredentials: make(map[string]bool),
	}

	verifier := NewVerifier(logger, publicKeyStore, statusChecker, revocationStore)

	// Verify presentation
	result, err := verifier.VerifyPresentation(presentation, nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Valid)
	assert.Empty(t, result.Errors)
	assert.True(t, result.Checks["basic_validation"])
	assert.True(t, result.Checks["presentation_proof"])
	assert.True(t, result.Checks["credentials_valid"])
}

func TestConstants(t *testing.T) {
	// Test default contexts
	contexts := DefaultContexts()
	assert.Contains(t, contexts, W3CCredentialsContext)
	assert.Contains(t, contexts, W3CSecurityContext)

	// Test default VC types
	vcTypes := DefaultVCTypes()
	assert.Contains(t, vcTypes, VerifiableCredentialType)

	// Test default VP types
	vpTypes := DefaultVPTypes()
	assert.Contains(t, vpTypes, VerifiablePresentationType)
}