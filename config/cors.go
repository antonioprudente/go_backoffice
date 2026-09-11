package config

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	// 1. Leggi gli origini permessi dall'ambiente
	rawOrigins := os.Getenv("ALLOWED_ORIGINS")

	var allowedOrigins []string
	if rawOrigins != "" {
		allowedOrigins = strings.Split(rawOrigins, ",")
	} else {
		// Fallback predefinito se la variabile d'ambiente non è impostata
		allowedOrigins = []string{"http://localhost:5173"}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 2. Verifica se l'origine della richiesta è tra quelli consentiti
		if origin != "" {
			for _, allowed := range allowedOrigins {
				allowed = strings.TrimSpace(allowed) // Rimuove eventuali spazi bianchi
				if origin == allowed || strings.HasSuffix(origin, ".netlify.app") {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		// 3. Gestione corretta Preflight OPTIONS
		if c.Request.Method == http.MethodOptions {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}

		c.Next()
	}
}
