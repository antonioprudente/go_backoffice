package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS gestisce le intestazioni CORS per il frontend React
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigin := "http://localhost:5174" // Sostituisci con l'URL esatto del tuo frontend

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		// Gestione della richiesta Preflight OPTIONS
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent) // 204 No Content
			return
		}

		c.Next()
	}
}
