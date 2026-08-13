package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRoles consente l'accesso solo se userRole (impostato da AuthMiddleware)
// è tra i ruoli passati come parametro
func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(allowedRoles))
	for _, r := range allowedRoles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		roleVal, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Utente non autenticato"})
			return
		}

		role, ok := roleVal.(string)
		if !ok || !allowed[role] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Accesso negato per il tuo ruolo"})
			return
		}

		c.Next()
	}
}
