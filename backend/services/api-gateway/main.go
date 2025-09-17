package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"

	m "github.com/verza-platform/verza/services/api-gateway/internal/models"
)

type Config struct {
	Port string `envconfig:"PORT" default:"8080"`
	Env  string `envconfig:"ENV" default:"dev"`
}

func main() {
	var cfg Config
	_ = envconfig.Process("GATEWAY", &cfg)

	logger, _ := zap.NewProduction()
	if cfg.Env == "dev" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		c.Next()
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/v1")
	{
		v1.POST("/kyc/submit", kycSubmit)
		v1.GET("/kyc/result/:jobId", kycResult)
		v1.POST("/vc/issue", vcIssue)
		v1.POST("/vp/verify", vpVerify)
		v1.POST("/vc/revoke", vcRevoke)
		v1.GET("/status/:vcHash", statusGet)
		v1.POST("/did/resolve", didResolve)
		v1.POST("/backup/register", backupRegister)
		v1.POST("/auth/did-challenge", didChallenge)
		v1.POST("/auth/did-response", didResponse)
	}

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	// Graceful shutdown placeholder
	done := make(chan os.Signal, 1)
	<-done
	_ = srv.Shutdown(context.Background())
}

// Handlers (stubs returning shapes that match the prompt)

func kycSubmit(c *gin.Context) {
	var req m.KYCSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m.KYCSubmitResponse{JobID: "job_mock_123"})
}

func kycResult(c *gin.Context) {
	resp := m.KYCResultResponse{Score: 0.98, Liveness: true, DocValid: true, OCR: map[string]string{"name": "Alice"}}
	c.JSON(http.StatusOK, resp)
}

func vcIssue(c *gin.Context) {
	var req m.IssueVCRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	vc := map[string]interface{}{"@context": []string{"https://www.w3.org/2018/credentials/v1"}, "type": []string{"VerifiableCredential", "KYCVerification"}, "issuer": "did:web:onboardia.com", "credentialSubject": map[string]interface{}{"id": req.SubjectDID, "claims": req.Claims}, "expirationDate": time.Now().Add(365*24*time.Hour).UTC().Format(time.RFC3339)}
	c.JSON(http.StatusOK, m.IssueVCResponse{VC: vc, AnchorTx: m.AnchorTx{ChainID: "polygon-mumbai", TxHash: "0xmock"}})
}

func vpVerify(c *gin.Context) {
	var req m.VPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, m.VPVerifyResponse{Valid: true, Details: map[string]interface{}{"note": "stub"}})
}

func vcRevoke(c *gin.Context) {
	var req m.RevokeVCRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	now := time.Now().UTC()
	c.JSON(http.StatusOK, m.StatusResponse{Anchored: true, Revoked: true, IssuedAt: now.Add(-time.Hour), RevokedAt: &now, AnchorTx: &m.AnchorTx{ChainID: "polygon-mumbai", TxHash: "0xrevoke"}, URI: "ipfs://statuslist"})
}

func statusGet(c *gin.Context) {
	now := time.Now().UTC()
	c.JSON(http.StatusOK, m.StatusResponse{Anchored: true, Revoked: false, IssuedAt: now.Add(-2 * time.Hour), AnchorTx: &m.AnchorTx{ChainID: "polygon-mumbai", TxHash: "0xanchor"}, URI: "ipfs://statuslist"})
}

func didResolve(c *gin.Context) {
	var req m.DIDResolveRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	doc := m.DIDDocument{ID: req.Did, VerificationMethod: []map[string]interface{}{{"id": req.Did + "#keys-1", "type": "Ed25519VerificationKey2020"}}}
	c.JSON(http.StatusOK, doc)
}

func backupRegister(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) }

func didChallenge(c *gin.Context) {
	var req m.DIDChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, m.DIDChallengeResponse{Challenge: "nonce-123", Domain: "verza.dev"})
}

func didResponse(c *gin.Context) {
	var req m.DIDResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, m.DIDResponseToken{Token: "jwt-mock"})
}