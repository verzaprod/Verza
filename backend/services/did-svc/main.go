package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/resolve", func(c *gin.Context) {
		var body struct{ DID string `json:"did" binding:"required"` }
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": body.DID, "verificationMethod": []map[string]string{{"id": body.DID + "#keys-1", "type": "Ed25519VerificationKey2020"}}})
	})
	r.Run(":8081")
}