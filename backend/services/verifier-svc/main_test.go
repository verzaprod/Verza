package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"github.com/verza/pkg/vc"
)

func init() {
	// Initialize logger for tests
	logger, _ = zap.NewDevelopment()
	
	// Initialize verifier for tests
	verifier = vc.NewVerifier(
		logger,
		&mockPublicKeyStore{},
		&mockStatusChecker{},
		&mockRevocationStore{},
	)
}

func TestHealthz(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/healthz", healthzHandler)

	req, _ := http.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestVerifyVP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/v1/vp/verify", verifyVPHandler)

	// Create a test credential
	testTime, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:00Z")
	testCredential := vc.VerifiableCredential{
		Context: []string{"https://www.w3.org/2018/credentials/v1"},
		Type:    []string{"VerifiableCredential"},
		Issuer:  vc.CredentialIssuer{ID: "did:example:issuer"},
		IssuanceDate: testTime,
		CredentialSubject: map[string]interface{}{
			"id": "did:example:subject",
			"name": "Test Subject",
		},
		Proof: &vc.Proof{
			Type: "JsonWebSignature2020",
			Created: testTime,
			VerificationMethod: "did:example:issuer#key-1",
			ProofPurpose: "assertionMethod",
			JWS: "test-signature",
		},
	}

	reqBody := VerifyVPRequest{
		VP: &vc.VerifiablePresentation{
			Context: []string{"https://www.w3.org/2018/credentials/v1"},
			Type:    []string{"VerifiablePresentation"},
			Holder:  "did:key:z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK",
			VerifiableCredential: []vc.VerifiableCredential{testCredential},
			Proof: &vc.Proof{
				Type: "JsonWebSignature2020",
				Created: testTime,
				VerificationMethod: "did:key:z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK#key-1",
				ProofPurpose: "authentication",
				Challenge: "test-challenge",
				Domain: "example.com",
				JWS: "test-vp-signature",
			},
		},
		Options: VerifyOptions{
			Challenge: "test-challenge",
			Domain:    "example.com",
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/v1/vp/verify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response VerifyVPResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Valid)
	assert.Equal(t, "VP verification successful (stub)", response.Details)
}