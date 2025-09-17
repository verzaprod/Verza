package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/verza/pkg/kms"
	"github.com/verza/pkg/vc"
	"go.uber.org/zap"
)

var (
	testIssuer *vc.KMSIssuer
)

func init() {
	// Initialize logger for tests
	logger, _ = zap.NewDevelopment()
	
	// Initialize KMS client for tests
	config := kms.Config{Provider: kms.ProviderLocal}
	kmsClient, err := kms.NewFactory().Create(logger, config)
	if err != nil {
		panic(err)
	}
	
	// Ensure test key exists
	keyID := "test-issuer-key"
	if _, err := kmsClient.GetKeyInfo(context.Background(), keyID); err != nil {
		if _, err := kmsClient.CreateKey(context.Background(), keyID, kms.KeyTypeRSA2048); err != nil {
			panic(err)
		}
	}
	
	// Initialize KMS issuer for tests
	testIssuer, err = vc.NewKMSIssuer(logger, "did:example:issuer", kmsClient, keyID, kms.AlgRS256)
	if err != nil {
		panic(err)
	}
	
	// Set global issuer for handlers
	vcIssuer = testIssuer
}

func TestHealthz(t *testing.T){
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/healthz", func(c *gin.Context){ c.JSON(200, gin.H{"status":"ok"}) })

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 { t.Fatalf("expected 200, got %d", w.Code) }
}

func TestIssueStub(t *testing.T){
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/v1/vc/issue", issueHandler)
	body := bytes.NewBufferString(`{"subjectDID":"did:key:abc","claims":{"kycLevel":"basic"}}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/vc/issue", body)
	req.Header.Set("Content-Type","application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 { t.Fatalf("expected 200, got %d", w.Code) }
}

func TestRevokeStub(t *testing.T){
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/v1/vc/revoke", revokeHandler)
	body := bytes.NewBufferString(`{"vcHash":"0xdeadbeef"}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/vc/revoke", body)
	req.Header.Set("Content-Type","application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 { t.Fatalf("expected 200, got %d", w.Code) }
}