package vc

import (
	"encoding/json"
	"time"
)

// VerifiableCredential represents a W3C Verifiable Credential
type VerifiableCredential struct {
	Context           []string               `json:"@context"`
	ID                string                 `json:"id,omitempty"`
	Type              []string               `json:"type"`
	Issuer            CredentialIssuer       `json:"issuer"`
	IssuanceDate      time.Time              `json:"issuanceDate"`
	ExpirationDate    *time.Time             `json:"expirationDate,omitempty"`
	CredentialSubject interface{}            `json:"credentialSubject"`
	CredentialStatus  *CredentialStatus      `json:"credentialStatus,omitempty"`
	Proof             *Proof                 `json:"proof,omitempty"`
	RefreshService    *RefreshService        `json:"refreshService,omitempty"`
	TermsOfUse        []TermsOfUse           `json:"termsOfUse,omitempty"`
	Evidence          []Evidence             `json:"evidence,omitempty"`
}

// VerifiablePresentation represents a W3C Verifiable Presentation
type VerifiablePresentation struct {
	Context              []string                `json:"@context"`
	ID                   string                  `json:"id,omitempty"`
	Type                 []string                `json:"type"`
	Holder               string                  `json:"holder,omitempty"`
	VerifiableCredential []VerifiableCredential `json:"verifiableCredential,omitempty"`
	Proof                *Proof                  `json:"proof,omitempty"`
}

// CredentialIssuer represents the issuer of a credential
type CredentialIssuer struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// CredentialSubject represents the subject of a credential
type CredentialSubject struct {
	ID   string                 `json:"id,omitempty"`
	Data map[string]interface{} `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for CredentialSubject
func (cs CredentialSubject) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	if cs.ID != "" {
		result["id"] = cs.ID
	}
	for k, v := range cs.Data {
		result[k] = v
	}
	return json.Marshal(result)
}

// UnmarshalJSON implements custom JSON unmarshaling for CredentialSubject
func (cs *CredentialSubject) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if id, ok := raw["id"].(string); ok {
		cs.ID = id
		delete(raw, "id")
	}

	cs.Data = raw
	return nil
}

// CredentialStatus represents the status information for a credential
type CredentialStatus struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// Proof represents a cryptographic proof
type Proof struct {
	Type               string    `json:"type"`
	Created            time.Time `json:"created"`
	VerificationMethod string    `json:"verificationMethod"`
	ProofPurpose       string    `json:"proofPurpose"`
	JWS                string    `json:"jws,omitempty"`
	ProofValue         string    `json:"proofValue,omitempty"`
	Challenge          string    `json:"challenge,omitempty"`
	Domain             string    `json:"domain,omitempty"`
}

// RefreshService represents a refresh service for a credential
type RefreshService struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// TermsOfUse represents terms of use for a credential
type TermsOfUse struct {
	Type string `json:"type"`
	ID   string `json:"id,omitempty"`
}

// Evidence represents evidence for a credential
type Evidence struct {
	Type string                 `json:"type"`
	ID   string                 `json:"id,omitempty"`
	Data map[string]interface{} `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for Evidence
func (e Evidence) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	result["type"] = e.Type
	if e.ID != "" {
		result["id"] = e.ID
	}
	for k, v := range e.Data {
		result[k] = v
	}
	return json.Marshal(result)
}

// UnmarshalJSON implements custom JSON unmarshaling for Evidence
func (e *Evidence) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if typ, ok := raw["type"].(string); ok {
		e.Type = typ
		delete(raw, "type")
	}

	if id, ok := raw["id"].(string); ok {
		e.ID = id
		delete(raw, "id")
	}

	e.Data = raw
	return nil
}

// VerificationOptions represents options for verification
type VerificationOptions struct {
	Challenge          string    `json:"challenge,omitempty"`
	Domain             string    `json:"domain,omitempty"`
	VerificationMethod string    `json:"verificationMethod,omitempty"`
	Created            time.Time `json:"created,omitempty"`
}

// SigningOptions represents options for signing
type SigningOptions struct {
	VerificationMethod string    `json:"verificationMethod"`
	ProofPurpose       string    `json:"proofPurpose"`
	Challenge          string    `json:"challenge,omitempty"`
	Domain             string    `json:"domain,omitempty"`
	Created            time.Time `json:"created,omitempty"`
}

// Constants for W3C VC contexts and types
const (
	// Standard W3C contexts
	W3CCredentialsContext = "https://www.w3.org/2018/credentials/v1"
	W3CSecurityContext    = "https://w3id.org/security/v1"
	W3CDIDContext         = "https://www.w3.org/ns/did/v1"

	// Standard credential types
	VerifiableCredentialType   = "VerifiableCredential"
	VerifiablePresentationType = "VerifiablePresentation"

	// Proof types
	JSONWebSignature2020 = "JsonWebSignature2020"
	Ed25519Signature2018 = "Ed25519Signature2018"
	RsaSignature2018     = "RsaSignature2018"

	// Proof purposes
	AssertionMethod      = "assertionMethod"
	Authentication       = "authentication"
	KeyAgreement         = "keyAgreement"
	CapabilityInvocation = "capabilityInvocation"
	CapabilityDelegation = "capabilityDelegation"

	// Status types
	RevocationList2020Status = "RevocationList2020Status"
	StatusList2021Entry      = "StatusList2021Entry"
)

// DefaultContexts returns the default contexts for VCs
func DefaultContexts() []string {
	return []string{
		W3CCredentialsContext,
		W3CSecurityContext,
	}
}

// DefaultVCTypes returns the default types for VCs
func DefaultVCTypes() []string {
	return []string{
		VerifiableCredentialType,
	}
}

// DefaultVPTypes returns the default types for VPs
func DefaultVPTypes() []string {
	return []string{
		VerifiablePresentationType,
	}
}