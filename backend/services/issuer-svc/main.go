package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"

	"github.com/verza/pkg/kms"
	"github.com/verza/pkg/vc"
)

type Config struct {
	Port      string `envconfig:"PORT" default:"8082"`
	Env       string `envconfig:"ENV" default:"dev"`
	IssuerDID string `envconfig:"ISSUER_DID" default:"did:web:verza.io:issuer"`
	KeyID     string `envconfig:"KEY_ID" default:"issuer-key-1"`
	KMS       kms.Config
}

var (
	logger    *zap.Logger
	vcIssuer  *vc.KMSIssuer
	kmsClient kms.KMS
)

// Request/response shapes (minimal) align with gateway models
 type IssueVCRequest struct {
 	SubjectDID string                 `json:"subjectDID" binding:"required"`
 	Claims     map[string]interface{} `json:"claims" binding:"required"`
 	Expiry     *string                `json:"expiry"`
 }
 type AnchorTx struct { ChainID string `json:"chainId"`; TxHash string `json:"txHash"` }
 type IssueVCResponse struct { VC *vc.VerifiableCredential `json:"vc"`; AnchorTx AnchorTx `json:"anchorTx"` }
 type RevokeVCRequest struct { VCHash string `json:"vcHash" binding:"required"`; Reason *string `json:"reason"` }
 type StatusResponse struct {
 	Anchored bool      `json:"anchored"`
 	Revoked  bool      `json:"revoked"`
 	IssuedAt time.Time `json:"issuedAt"`
 	RevokedAt *time.Time `json:"revokedAt,omitempty"`
 	AnchorTx *AnchorTx `json:"anchorTx,omitempty"`
 	URI      string    `json:"uri"`
 }

func main() {
	var cfg Config
	_ = envconfig.Process("ISSUER_SVC", &cfg)

	logger, _ = zap.NewDevelopment()
	if cfg.Env != "dev" { logger, _ = zap.NewProduction() }
	defer logger.Sync()

	// Initialize KMS
	kmsFactory := kms.NewFactory()
	var err error
	kmsClient, err = kmsFactory.Create(logger, cfg.KMS)
	if err != nil {
		logger.Fatal("Failed to create KMS client", zap.Error(err))
	}

	// Ensure issuer key exists in KMS
	ctx := context.Background()
	_, err = kmsClient.GetKeyInfo(ctx, cfg.KeyID)
	if err != nil {
		logger.Info("Creating new issuer key in KMS", zap.String("key_id", cfg.KeyID))
		_, err = kmsClient.CreateKey(ctx, cfg.KeyID, kms.KeyTypeRSA2048)
		if err != nil {
			logger.Fatal("Failed to create issuer key in KMS", zap.Error(err))
		}
	}

	// Initialize KMS-based VC issuer
	vcIssuer, err = vc.NewKMSIssuer(logger, cfg.IssuerDID, kmsClient, cfg.KeyID, kms.AlgRS256)
	if err != nil {
		logger.Fatal("Failed to create KMS VC issuer", zap.Error(err))
	}

	logger.Info("KMS VC Issuer initialized", 
		zap.String("issuerDID", cfg.IssuerDID),
		zap.String("keyID", cfg.KeyID),
		zap.String("kms_provider", string(cfg.KMS.Provider)))

	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/healthz", healthzHandler)

	v1 := r.Group("/v1")
	{
		vcGroup := v1.Group("/vc")
		vcGroup.POST("/issue", issueHandler)
		vcGroup.POST("/revoke", revokeHandler)
		vcGroup.GET("/status/:id", statusHandler)
	}

	logger.Info("Starting issuer service", zap.String("port", cfg.Port))
	_ = r.Run(":" + cfg.Port)
}

func healthzHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "issuer-svc"})
}

func issueHandler(c *gin.Context) {
	var req IssueVCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Issuing credential with KMS", 
		zap.String("subjectDID", req.SubjectDID),
		zap.Any("claims", req.Claims))

	// Create credential subject
	subject := vc.CredentialSubject{
		ID:   req.SubjectDID,
		Data: req.Claims,
	}

	// Determine credential type based on claims
	credentialTypes := []string{"IdentityCredential"}
	if _, hasKYC := req.Claims["kycVerified"]; hasKYC {
		credentialTypes = []string{"KYCCredential"}
	}

	// Create issue request
	issueReq := vc.IssueCredentialRequest{
		SubjectDID:        req.SubjectDID,
		CredentialTypes:   credentialTypes,
		CredentialSubject: &subject,
	}

	// Set expiration if provided
	if req.Expiry != nil {
		if expiryTime, err := time.Parse(time.RFC3339, *req.Expiry); err == nil {
			issueReq.ExpirationDate = &expiryTime
		} else {
			logger.Warn("Invalid expiry format, ignoring", zap.String("expiry", *req.Expiry))
		}
	}

	// Issue the credential using KMS
	ctx := context.WithValue(c.Request.Context(), "request_id", c.GetHeader("X-Request-ID"))
	credential, err := vcIssuer.IssueCredential(ctx, issueReq)
	if err != nil {
		logger.Error("Failed to issue credential", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to issue credential"})
		return
	}

	// Mock blockchain anchoring (would be real in production)
	anchorTx := AnchorTx{
		ChainID: "polygon-mumbai",
		TxHash:  fmt.Sprintf("0x%x", time.Now().UnixNano()),
	}

	resp := IssueVCResponse{
		VC:       credential,
		AnchorTx: anchorTx,
	}

	logger.Info("Credential issued successfully", 
		zap.String("credentialID", credential.ID),
		zap.String("txHash", anchorTx.TxHash))

	c.JSON(http.StatusOK, resp)
}

func revokeHandler(c *gin.Context) {
	var req RevokeVCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid revoke request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info("Revoking credential", 
		zap.String("vcHash", req.VCHash),
		zap.String("reason", func() string {
			if req.Reason != nil {
				return *req.Reason
			}
			return "No reason provided"
		}()))

	// Mock revocation (would interact with blockchain in production)
	now := time.Now().UTC()
	issuedAt := now.Add(-time.Hour) // Mock issued time

	// Mock blockchain transaction for revocation
	anchorTx := &AnchorTx{
		ChainID: "polygon-mumbai",
		TxHash:  fmt.Sprintf("0x%x", time.Now().UnixNano()),
	}

	resp := StatusResponse{
		Anchored:  true,
		Revoked:   true,
		IssuedAt:  issuedAt,
		RevokedAt: &now,
		AnchorTx:  anchorTx,
		URI:       fmt.Sprintf("https://verza.io/status/%s", req.VCHash),
	}

	logger.Info("Credential revoked successfully", 
		zap.String("vcHash", req.VCHash),
		zap.String("txHash", anchorTx.TxHash))

	c.JSON(http.StatusOK, resp)
}

func statusHandler(c *gin.Context) {
	credentialID := c.Param("id")
	if credentialID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Credential ID is required"})
		return
	}

	logger.Info("Checking credential status", zap.String("credentialID", credentialID))

	// Mock status check (would query blockchain/database in production)
	issuedAt := time.Now().Add(-2 * time.Hour)
	anchorTx := &AnchorTx{
		ChainID: "polygon-mumbai",
		TxHash:  "0xmockstatus",
	}

	resp := StatusResponse{
		Anchored:  true,
		Revoked:   false,
		IssuedAt:  issuedAt,
		RevokedAt: nil,
		AnchorTx:  anchorTx,
		URI:       fmt.Sprintf("https://verza.io/status/%s", credentialID),
	}

	c.JSON(http.StatusOK, resp)
}