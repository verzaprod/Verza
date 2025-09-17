package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Config struct {
	Port        string `envconfig:"PORT" default:"8084"`
	EthRPCURL   string `envconfig:"ETH_RPC_URL" default:"http://localhost:8545"`
	ContractAddr string `envconfig:"CONTRACT_ADDR" default:"0x0000000000000000000000000000000000000000"`
	PrivateKey  string `envconfig:"PRIVATE_KEY" default:""`
}

type AnchorRequest struct {
	VCHash string `json:"vcHash" binding:"required"`
	URI    string `json:"uri" binding:"required"`
}

type AnchorResponse struct {
	TxHash  string `json:"txHash"`
	ChainID int64  `json:"chainId"`
	Status  string `json:"status"`
}

type RevokeRequest struct {
	VCHash string `json:"vcHash" binding:"required"`
	Reason string `json:"reason,omitempty"`
}

type RevokeResponse struct {
	TxHash  string `json:"txHash"`
	ChainID int64  `json:"chainId"`
	Status  string `json:"status"`
}

var (
	logger *zap.Logger
	config Config
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

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Health check endpoint
	router.GET("/healthz", healthzHandler)

	// Anchor service endpoints
	v1 := router.Group("/v1")
	{
		v1.POST("/anchor", anchorHandler)
		v1.POST("/revoke", revokeHandler)
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	go func() {
		logger.Info("Starting anchor service", zap.String("port", config.Port))
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
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "anchor-svc"})
}

func anchorHandler(c *gin.Context) {
	var req AnchorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Anchoring VC", zap.String("vcHash", req.VCHash), zap.String("uri", req.URI))

	// TODO: Implement actual blockchain anchoring
	// For now, return a mock response
	response := AnchorResponse{
		TxHash:  "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		ChainID: 137, // Polygon mainnet
		Status:  "anchored",
	}

	c.JSON(http.StatusOK, response)
}

func revokeHandler(c *gin.Context) {
	var req RevokeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Revoking VC", zap.String("vcHash", req.VCHash), zap.String("reason", req.Reason))

	// TODO: Implement actual blockchain revocation
	// For now, return a mock response
	response := RevokeResponse{
		TxHash:  "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		ChainID: 137, // Polygon mainnet
		Status:  "revoked",
	}

	c.JSON(http.StatusOK, response)
}