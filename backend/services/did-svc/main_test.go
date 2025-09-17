package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResolve(t *testing.T){
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/resolve", func(c *gin.Context) {
		var body struct{ DID string `json:"did" binding:"required"` }
		if err := c.ShouldBindJSON(&body); err != nil { c.JSON(400, gin.H{"error":err.Error()}); return }
		c.JSON(200, gin.H{"id": body.DID})
	})

	payload := bytes.NewBufferString(`{"did":"did:key:xyz"}`)
	req := httptest.NewRequest(http.MethodPost, "/resolve", payload)
	req.Header.Set("Content-Type","application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 { t.Fatalf("expected 200, got %d", w.Code) }
}