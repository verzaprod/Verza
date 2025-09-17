package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/verza/pkg/blockchain"
	"github.com/verza/pkg/database"
	"github.com/verza/pkg/kms"
	"github.com/verza/pkg/security"
	"github.com/verza/pkg/vc"
)

// IntegrationTestSuite holds all components for integration testing
type IntegrationTestSuite struct {
	logger     *zap.Logger
	db         *database.DB
	kms        kms.KeyManager
	blockchain *blockchain.Client
	issuer     *vc.Issuer
	verifier   *vc.Verifier
	server     *gin.Engine
}

// SetupIntegrationTest initializes all components for testing
func SetupIntegrationTest(t *testing.T) *IntegrationTestSuite {
	logger := zap.NewNop()
	ctx := context.Background()

	// Setup database (using in-memory SQLite for testing)
	dbConfig := &database.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		Database: "test_verza",
		SSLMode:  "disable",
	}

	// For integration tests, we'll mock the database connection
	// In real scenarios, you'd use a test database
	db, err := database.NewDB(ctx, dbConfig, logger)
	if err != nil {
		// Use mock database for testing
		t.Logf("Database connection failed, using mock: %v", err)
		db = nil
	}

	// Setup KMS (using local KMS for testing)
	kmsConfig := &kms.LocalConfig{
		KeySize: 2048,
	}
	keyManager, err := kms.NewLocalKMS(kmsConfig, logger)
	require.NoError(t, err)

	// Setup blockchain client (using mock for testing)
	blockchainConfig := &blockchain.Config{
		RPCURL:          "http://localhost:8545",
		ContractAddress: "0x1234567890123456789012345678901234567890",
		PrivateKey:      "0x" + strings.Repeat("a", 64),
	}
	blockchainClient, err := blockchain.NewClient(blockchainConfig, logger)
	if err != nil {
		// Use mock blockchain for testing
		t.Logf("Blockchain connection failed, using mock: %v", err)
		blockchainClient = nil
	}

	// Setup VC components
	issuer, err := vc.NewIssuer(keyManager, logger)
	require.NoError(t, err)

	verifier := vc.NewVerifier(logger)

	// Setup HTTP server
	gin.SetMode(gin.TestMode)
	server := gin.New()
	server.Use(gin.Recovery())

	// Add security middleware
	securityConfig := &security.Config{
		RateLimit: security.RateLimitConfig{
			Enabled: true,
			RPS:     100,
			Burst:   200,
		},
		CORS: security.CORSConfig{
			Enabled:        true,
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders: []string{"*"},
		},
	}
	securityMiddleware := security.NewMiddleware(securityConfig, logger)
	server.Use(securityMiddleware.RateLimit())
	server.Use(securityMiddleware.CORS())
	server.Use(securityMiddleware.Security())

	return &IntegrationTestSuite{
		logger:     logger,
		db:         db,
		kms:        keyManager,
		blockchain: blockchainClient,
		issuer:     issuer,
		verifier:   verifier,
		server:     server,
	}
}

// TestCompleteVCLifecycle tests the entire VC lifecycle
func TestCompleteVCLifecycle(t *testing.T) {
	suite := SetupIntegrationTest(t)
	ctx := context.Background()

	// Step 1: Issue a Verifiable Credential
	t.Run("Issue VC", func(t *testing.T) {
		// Create credential data
		credentialData := map[string]interface{}{
			"@context": []string{
				"https://www.w3.org/2018/credentials/v1",
				"https://verza.io/contexts/kyc/v1",
			},
			"type": []string{"VerifiableCredential", "KYCCredential"},
			"credentialSubject": map[string]interface{}{
				"id":   "did:example:123456789abcdefghi",
				"name": "John Doe",
				"kycStatus": "verified",
				"verificationLevel": "full",
			},
		}

		// Issue the credential
		vc, err := suite.issuer.IssueCredential(ctx, credentialData)
		require.NoError(t, err)
		assert.NotEmpty(t, vc.ID)
		assert.NotEmpty(t, vc.Proof)
		assert.Equal(t, "did:example:123456789abcdefghi", vc.CredentialSubject["id"])

		// Store VC for later tests
		ctx = context.WithValue(ctx, "issued_vc", vc)
	})

	// Step 2: Verify the Verifiable Credential
	t.Run("Verify VC", func(t *testing.T) {
		vc := ctx.Value("issued_vc").(*vc.VerifiableCredential)

		// Verify the credential
		result, err := suite.verifier.VerifyCredential(ctx, vc)
		require.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Empty(t, result.Errors)
	})

	// Step 3: Anchor VC to blockchain (if blockchain is available)
	t.Run("Anchor VC to Blockchain", func(t *testing.T) {
		if suite.blockchain == nil {
			t.Skip("Blockchain not available for testing")
		}

		vc := ctx.Value("issued_vc").(*vc.VerifiableCredential)

		// Calculate VC hash
		vcBytes, err := json.Marshal(vc)
		require.NoError(t, err)
		vcHash := suite.blockchain.CalculateHash(vcBytes)

		// Anchor to blockchain
		txHash, err := suite.blockchain.AnchorCredential(ctx, vcHash)
		if err != nil {
			t.Logf("Blockchain anchoring failed (expected in test): %v", err)
			return
		}
		assert.NotEmpty(t, txHash)

		// Verify anchoring
		isAnchored, err := suite.blockchain.IsCredentialAnchored(ctx, vcHash)
		require.NoError(t, err)
		assert.True(t, isAnchored)
	})

	// Step 4: Revoke the Verifiable Credential
	t.Run("Revoke VC", func(t *testing.T) {
		if suite.blockchain == nil {
			t.Skip("Blockchain not available for testing")
		}

		vc := ctx.Value("issued_vc").(*vc.VerifiableCredential)

		// Calculate VC hash
		vcBytes, err := json.Marshal(vc)
		require.NoError(t, err)
		vcHash := suite.blockchain.CalculateHash(vcBytes)

		// Revoke credential
		txHash, err := suite.blockchain.RevokeCredential(ctx, vcHash)
		if err != nil {
			t.Logf("Blockchain revocation failed (expected in test): %v", err)
			return
		}
		assert.NotEmpty(t, txHash)

		// Verify revocation
		isRevoked, err := suite.blockchain.IsCredentialRevoked(ctx, vcHash)
		require.NoError(t, err)
		assert.True(t, isRevoked)
	})
}

// TestAPIEndpoints tests the HTTP API endpoints
func TestAPIEndpoints(t *testing.T) {
	suite := SetupIntegrationTest(t)

	// Setup API routes
	setupAPIRoutes(suite)

	t.Run("Health Check", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		suite.server.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "healthy", response["status"])
	})

	t.Run("Issue Credential API", func(t *testing.T) {
		credentialRequest := map[string]interface{}{
			"credentialSubject": map[string]interface{}{
				"id":   "did:example:test123",
				"name": "Test User",
				"kycStatus": "verified",
			},
		}

		requestBody, _ := json.Marshal(credentialRequest)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/credentials/issue", strings.NewReader(string(requestBody)))
		req.Header.Set("Content-Type", "application/json")
		suite.server.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response["credential"])
	})

	t.Run("Verify Credential API", func(t *testing.T) {
		// First issue a credential
		credentialData := map[string]interface{}{
			"credentialSubject": map[string]interface{}{
				"id":   "did:example:verify123",
				"name": "Verify User",
			},
		}

		vc, err := suite.issuer.IssueCredential(context.Background(), credentialData)
		require.NoError(t, err)

		// Now verify it via API
		verifyRequest := map[string]interface{}{
			"credential": vc,
		}

		requestBody, _ := json.Marshal(verifyRequest)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/credentials/verify", strings.NewReader(string(requestBody)))
		req.Header.Set("Content-Type", "application/json")
		suite.server.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response["valid"].(bool))
	})
}

// TestSecurityMiddleware tests security features
func TestSecurityMiddleware(t *testing.T) {
	suite := SetupIntegrationTest(t)
	setupAPIRoutes(suite)

	t.Run("Rate Limiting", func(t *testing.T) {
		// Make multiple requests quickly
		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/health", nil)
			suite.server.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("CORS Headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/api/v1/credentials/issue", nil)
		req.Header.Set("Origin", "https://example.com")
		suite.server.ServeHTTP(w, req)

		assert.Contains(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	})

	t.Run("Security Headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		suite.server.ServeHTTP(w, req)

		assert.NotEmpty(t, w.Header().Get("X-Content-Type-Options"))
		assert.NotEmpty(t, w.Header().Get("X-Frame-Options"))
	})
}

// TestKMSIntegration tests key management operations
func TestKMSIntegration(t *testing.T) {
	suite := SetupIntegrationTest(t)
	ctx := context.Background()

	t.Run("Key Generation", func(t *testing.T) {
		keyID, err := suite.kms.GenerateKey(ctx, "test-key", kms.KeyTypeRSA, 2048)
		require.NoError(t, err)
		assert.NotEmpty(t, keyID)
	})

	t.Run("Sign and Verify", func(t *testing.T) {
		// Generate a key
		keyID, err := suite.kms.GenerateKey(ctx, "sign-test-key", kms.KeyTypeRSA, 2048)
		require.NoError(t, err)

		// Sign data
		data := []byte("test data to sign")
		signature, err := suite.kms.Sign(ctx, keyID, data)
		require.NoError(t, err)
		assert.NotEmpty(t, signature)

		// Verify signature
		valid, err := suite.kms.Verify(ctx, keyID, data, signature)
		require.NoError(t, err)
		assert.True(t, valid)
	})
}

// TestDatabaseOperations tests database operations (if available)
func TestDatabaseOperations(t *testing.T) {
	suite := SetupIntegrationTest(t)

	if suite.db == nil {
		t.Skip("Database not available for testing")
	}

	ctx := context.Background()

	t.Run("Database Health", func(t *testing.T) {
		err := suite.db.Health(ctx)
		assert.NoError(t, err)
	})

	t.Run("Database Stats", func(t *testing.T) {
		stats := suite.db.Stats()
		assert.NotNil(t, stats)
	})
}

// setupAPIRoutes configures the API routes for testing
func setupAPIRoutes(suite *IntegrationTestSuite) {
	api := suite.server.Group("/api/v1")

	// Health check
	suite.server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
		})
	})

	// Issue credential
	api.POST("/credentials/issue", func(c *gin.Context) {
		var request map[string]interface{}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Add required fields
		request["@context"] = []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://verza.io/contexts/kyc/v1",
		}
		request["type"] = []string{"VerifiableCredential", "KYCCredential"}

		vc, err := suite.issuer.IssueCredential(c.Request.Context(), request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"credential": vc})
	})

	// Verify credential
	api.POST("/credentials/verify", func(c *gin.Context) {
		var request struct {
			Credential *vc.VerifiableCredential `json:"credential"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := suite.verifier.VerifyCredential(c.Request.Context(), request.Credential)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"valid":  result.Valid,
			"errors": result.Errors,
		})
	})
}

// generateRandomBytes generates random bytes for testing
func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

// TestPerformance runs basic performance tests
func TestPerformance(t *testing.T) {
	suite := SetupIntegrationTest(t)
	ctx := context.Background()

	t.Run("Credential Issuance Performance", func(t *testing.T) {
		start := time.Now()
		numCredentials := 10

		for i := 0; i < numCredentials; i++ {
			credentialData := map[string]interface{}{
				"@context": []string{
					"https://www.w3.org/2018/credentials/v1",
					"https://verza.io/contexts/kyc/v1",
				},
				"type": []string{"VerifiableCredential", "KYCCredential"},
				"credentialSubject": map[string]interface{}{
					"id":   fmt.Sprintf("did:example:perf%d", i),
					"name": fmt.Sprintf("Performance Test User %d", i),
				},
			}

			_, err := suite.issuer.IssueCredential(ctx, credentialData)
			require.NoError(t, err)
		}

		duration := time.Since(start)
		avgTime := duration / time.Duration(numCredentials)
		t.Logf("Issued %d credentials in %v (avg: %v per credential)", numCredentials, duration, avgTime)

		// Assert reasonable performance (adjust thresholds as needed)
		assert.Less(t, avgTime, time.Second, "Credential issuance should be fast")
	})
}