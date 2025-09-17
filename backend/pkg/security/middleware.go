package security

import (
	"crypto/tls"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	golang_rate "golang.org/x/time/rate"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     golang_rate.Limit
	burst    int
	cleanup  time.Duration
}

// Visitor represents a rate-limited client
type Visitor struct {
	limiter  *golang_rate.Limiter
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     golang_rate.Limit(requestsPerSecond),
		burst:    burst,
		cleanup:  time.Minute * 3,
	}
	
	// Start cleanup goroutine
	go rl.cleanupVisitors()
	
	return rl
}

// getVisitor returns the rate limiter for a given IP
func (rl *RateLimiter) getVisitor(ip string) *golang_rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	v, exists := rl.visitors[ip]
	if !exists {
		limiter := golang_rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = &Visitor{limiter, time.Now()}
		return limiter
	}
	
	v.lastSeen = time.Now()
	return v.limiter
}

// cleanupVisitors removes old visitors
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.cleanup {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware returns a Gin middleware for rate limiting
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)
		limiter := rl.getVisitor(ip)
		
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// getClientIP extracts the real client IP
func getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	// Check X-Real-IP header
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}
	
	// Fallback to RemoteAddr
	return c.ClientIP()
}

// SecurityHeaders middleware adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Enforce HTTPS
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")
		
		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		c.Next()
	}
}

// InputValidation middleware validates request size and content type
func InputValidation(maxBodySize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Limit request body size
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodySize)
		
		// Validate content type for POST/PUT requests
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if contentType == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Content-Type header is required",
					"code":  "MISSING_CONTENT_TYPE",
				})
				c.Abort()
				return
			}
			
			// Allow JSON and multipart form data
			if !strings.Contains(contentType, "application/json") && 
			   !strings.Contains(contentType, "multipart/form-data") {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "Unsupported content type",
					"code":  "UNSUPPORTED_MEDIA_TYPE",
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}

// APIKeyAuth middleware for API key authentication
type APIKeyAuth struct {
	validKeys map[string]bool
	headerName string
	logger    *zap.Logger
}

// NewAPIKeyAuth creates a new API key authentication middleware
func NewAPIKeyAuth(keys []string, headerName string, logger *zap.Logger) *APIKeyAuth {
	validKeys := make(map[string]bool)
	for _, key := range keys {
		validKeys[key] = true
	}
	
	return &APIKeyAuth{
		validKeys:  validKeys,
		headerName: headerName,
		logger:     logger,
	}
}

// Middleware returns the API key authentication middleware
func (a *APIKeyAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(a.headerName)
		if apiKey == "" {
			a.logger.Warn("Missing API key", zap.String("ip", getClientIP(c)))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API key required",
				"code":  "MISSING_API_KEY",
			})
			c.Abort()
			return
		}
		
		if !a.validKeys[apiKey] {
			a.logger.Warn("Invalid API key", 
				zap.String("ip", getClientIP(c)),
				zap.String("key_prefix", apiKey[:min(len(apiKey), 8)]),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
				"code":  "INVALID_API_KEY",
			})
			c.Abort()
			return
		}
		
		// Store API key info in context for later use
		c.Set("api_key", apiKey)
		c.Next()
	}
}

// TLSConfig returns a secure TLS configuration
func TLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

// RequestLogger middleware for structured request logging
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("HTTP Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
			zap.Int("body_size", param.BodySize),
		)
		return ""
	})
}

// CORS middleware for Cross-Origin Resource Sharing
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		c.Header("Access-Control-Max-Age", "86400")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}