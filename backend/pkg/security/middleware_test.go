package security

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create rate limiter: 2 requests per second, burst of 1
	rl := NewRateLimiter(2.0, 1)
	
	router := gin.New()
	router.Use(rl.RateLimitMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// First request should succeed
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	
	if w1.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w1.Code)
	}
	
	// Second request immediately should be rate limited
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", w2.Code)
	}
	
	// Different IP should not be rate limited
	req3 := httptest.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "192.168.1.2:12345"
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	
	if w3.Code != http.StatusOK {
		t.Errorf("Expected status 200 for different IP, got %d", w3.Code)
	}
}

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.Use(SecurityHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Check security headers
	expectedHeaders := map[string]string{
		"X-Frame-Options":           "DENY",
		"X-Content-Type-Options":     "nosniff",
		"X-XSS-Protection":           "1; mode=block",
		"Strict-Transport-Security":  "max-age=31536000; includeSubDomains",
		"Referrer-Policy":            "strict-origin-when-cross-origin",
	}
	
	for header, expectedValue := range expectedHeaders {
		actualValue := w.Header().Get(header)
		if actualValue != expectedValue {
			t.Errorf("Expected %s header to be %s, got %s", header, expectedValue, actualValue)
		}
	}
}

func TestInputValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.Use(InputValidation(1024)) // 1KB limit
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Test missing content type
	req1 := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"test": "data"}`))
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	
	if w1.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing content type, got %d", w1.Code)
	}
	
	// Test unsupported content type
	req2 := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"test": "data"}`))
	req2.Header.Set("Content-Type", "text/plain")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	if w2.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Expected status 415 for unsupported content type, got %d", w2.Code)
	}
	
	// Test valid JSON content type
	req3 := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"test": "data"}`))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	
	if w3.Code != http.StatusOK {
		t.Errorf("Expected status 200 for valid JSON, got %d", w3.Code)
	}
}

func TestAPIKeyAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	
	validKeys := []string{"test-key-123", "another-key-456"}
	auth := NewAPIKeyAuth(validKeys, "X-API-Key", logger)
	
	router := gin.New()
	router.Use(auth.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Test missing API key
	req1 := httptest.NewRequest("GET", "/test", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	
	if w1.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for missing API key, got %d", w1.Code)
	}
	
	// Test invalid API key
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("X-API-Key", "invalid-key")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	if w2.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for invalid API key, got %d", w2.Code)
	}
	
	// Test valid API key
	req3 := httptest.NewRequest("GET", "/test", nil)
	req3.Header.Set("X-API-Key", "test-key-123")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	
	if w3.Code != http.StatusOK {
		t.Errorf("Expected status 200 for valid API key, got %d", w3.Code)
	}
	
	// Check if API key is stored in context
	router2 := gin.New()
	router2.Use(auth.Middleware())
	router2.GET("/test", func(c *gin.Context) {
		apiKey, exists := c.Get("api_key")
		if !exists {
			t.Error("API key not found in context")
		}
		if apiKey != "test-key-123" {
			t.Errorf("Expected API key 'test-key-123', got %v", apiKey)
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	req4 := httptest.NewRequest("GET", "/test", nil)
	req4.Header.Set("X-API-Key", "test-key-123")
	w4 := httptest.NewRecorder()
	router2.ServeHTTP(w4, req4)
}

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	allowedOrigins := []string{"https://example.com", "https://app.example.com"}
	
	router := gin.New()
	router.Use(CORS(allowedOrigins))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	router.OPTIONS("/test", func(c *gin.Context) {
		// This should be handled by CORS middleware
	})
	
	// Test allowed origin
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.Header.Set("Origin", "https://example.com")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	
	if w1.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("Expected CORS header for allowed origin")
	}
	
	// Test disallowed origin
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Origin", "https://malicious.com")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	if w2.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("Expected no CORS header for disallowed origin")
	}
	
	// Test OPTIONS request
	req3 := httptest.NewRequest("OPTIONS", "/test", nil)
	req3.Header.Set("Origin", "https://example.com")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	
	if w3.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 for OPTIONS request, got %d", w3.Code)
	}
}

func TestGetClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	tests := []struct {
		name           string
		xForwardedFor  string
		xRealIP        string
		remoteAddr     string
		expectedIP     string
	}{
		{
			name:           "X-Forwarded-For header",
			xForwardedFor:  "203.0.113.1, 198.51.100.1",
			remoteAddr:     "192.168.1.1:12345",
			expectedIP:     "203.0.113.1",
		},
		{
			name:       "X-Real-IP header",
			xRealIP:    "203.0.113.2",
			remoteAddr: "192.168.1.1:12345",
			expectedIP: "203.0.113.2",
		},
		{
			name:       "RemoteAddr fallback",
			remoteAddr: "192.168.1.1:12345",
			expectedIP: "192.168.1.1",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/test", func(c *gin.Context) {
				ip := getClientIP(c)
				c.JSON(http.StatusOK, gin.H{"ip": ip})
			})
			
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}
			req.RemoteAddr = tt.remoteAddr
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			
			if response["ip"] != tt.expectedIP {
				t.Errorf("Expected IP %s, got %s", tt.expectedIP, response["ip"])
			}
		})
	}
}

func TestTLSConfig(t *testing.T) {
	config := TLSConfig()
	
	if config.MinVersion != 0x0303 { // TLS 1.2
		t.Errorf("Expected minimum TLS version 1.2, got %x", config.MinVersion)
	}
	
	if !config.PreferServerCipherSuites {
		t.Error("Expected PreferServerCipherSuites to be true")
	}
	
	if len(config.CipherSuites) == 0 {
		t.Error("Expected cipher suites to be configured")
	}
	
	if len(config.CurvePreferences) == 0 {
		t.Error("Expected curve preferences to be configured")
	}
}

func BenchmarkRateLimiter(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	rl := NewRateLimiter(1000.0, 100) // High limits for benchmarking
	
	router := gin.New()
	router.Use(rl.RateLimitMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.1:12345"
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}

func BenchmarkSecurityHeaders(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.Use(SecurityHeaders())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}