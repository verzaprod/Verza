package main

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"github.com/verza/pkg/vc"
)

type Config struct {
	Port string `envconfig:"PORT" default:"8083"`
	Env  string `envconfig:"ENV" default:"dev"`
}

type VerifyVPRequest struct {
	VP      *vc.VerifiablePresentation `json:"vp" binding:"required"`
	Options VerifyOptions              `json:"options"`
}

type VerifyOptions struct {
	Challenge string `json:"challenge"`
	Domain    string `json:"domain"`
}

type VerifyVPResponse struct {
	Valid   bool   `json:"valid"`
	Details string `json:"details"`
}

var (
	logger   *zap.Logger
	config   Config
	verifier *vc.Verifier
)

func main() {
	// Initialize logger
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Load configuration
	err = envconfig.Process("", &config)
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Initialize verifier with mock key store
	verifier = vc.NewVerifier(
		logger,
		&mockPublicKeyStore{},
		&mockStatusChecker{},
		&mockRevocationStore{},
	)

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Health check endpoint
	router.GET("/healthz", healthzHandler)

	// Verifier service endpoints
	v1 := router.Group("/v1")
	{
		v1.POST("/vp/verify", verifyVPHandler)
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	go func() {
		logger.Info("Starting verifier service", zap.String("port", config.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}

func healthzHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "verifier-svc"})
}

func verifyVPHandler(c *gin.Context) {
	var req VerifyVPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Verifying VP", zap.String("holder", req.VP.Holder), zap.String("challenge", req.Options.Challenge))

	// For testing purposes, always return success
	// In production, this would perform actual cryptographic verification
	logger.Info("VP verification (stub mode)", zap.String("holder", req.VP.Holder))
	
	response := VerifyVPResponse{
		Valid:   true,
		Details: "VP verification successful (stub)",
	}

	c.JSON(http.StatusOK, response)
}

// Mock implementations for testing
type mockPublicKeyStore struct{}

func (m *mockPublicKeyStore) GetPublicKey(keyID string) (crypto.PublicKey, error) {
	// Mock RSA public key for testing
	pubKeyPEM := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEjWT2btf+FxKaI9elKS9vhfPAHxunQgG7/TISLWlnDIOo
qbmn9TiJpkKqXiRmCnl1jLdKrttperNakajUyoL2jjjL1lUHBOCXjzw7/4/YrSQD
SuFqAiZjkrBgHU+mvGQxwEEEYZWEXwbeHhzrVB8plwqnI0cVhMDwrfmisIXhJDcR
Lrz8/ohm0W59Vqtb5qHlimudU3MDcPFuuaOKmMy6AB6VSWIFrPgNON8v69oQUNP
RCFN47SLlmx/VWjZlR1V8iNdX4pEVBzlEHfWaBw0HRWM7JX/xqVG7/tybg8hnH
GwIDAQAB
-----END PUBLIC KEY-----`
	block, _ := pem.Decode([]byte(pubKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	return pubKey.(*rsa.PublicKey), nil
}

type mockStatusChecker struct{}

func (m *mockStatusChecker) CheckStatus(statusListCredential string) (*vc.CredentialStatusInfo, error) {
	return &vc.CredentialStatusInfo{
		Valid: true,
	}, nil // Always return valid for mock
}

type mockRevocationStore struct{}

func (m *mockRevocationStore) IsRevoked(vcHash string) (bool, error) {
	return false, nil // Never revoked for mock
}