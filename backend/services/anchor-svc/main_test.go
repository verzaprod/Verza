package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func init() {
	// Initialize logger for tests
	logger, _ = zap.NewDevelopment()
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
	assert.Contains(t, w.Body.String(), "anchor-svc")
}

func TestAnchorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/v1/anchor", anchorHandler)

	reqBody := AnchorRequest{
		VCHash: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		URI:    "https://example.com/vc/status/123",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/v1/anchor", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response AnchorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.TxHash)
	assert.Equal(t, int64(137), response.ChainID)
	assert.Equal(t, "anchored", response.Status)
}

func TestRevokeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/v1/revoke", revokeHandler)

	reqBody := RevokeRequest{
		VCHash: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		Reason: "User requested revocation",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/v1/revoke", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response RevokeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.TxHash)
	assert.Equal(t, int64(137), response.ChainID)
	assert.Equal(t, "revoked", response.Status)
}

func TestAnchorHandlerBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/v1/anchor", anchorHandler)

	// Send invalid JSON
	req, _ := http.NewRequest("POST", "/v1/anchor", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRevokeHandlerBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/v1/revoke", revokeHandler)

	// Send invalid JSON
	req, _ := http.NewRequest("POST", "/v1/revoke", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}