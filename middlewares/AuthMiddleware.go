package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	secret := []byte(getEnv("JWT_SECRET", "super-secret-key"))

	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token mancante o formato non valido"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// 1. Controllo di sicurezza sull'algoritmo di firma
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metodo di firma non valido: %v", t.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token non valido o scaduto"})
			return
		}

		// 2. Estrazione sicura dei Claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato claims non valido"})
			return
		}

		// 3. Salvataggio di userID e userRole nel contesto di Gin
		if userID, exists := claims["sub"]; exists {
			ctx.Set("userID", userID)
		}

		if role, exists := claims["role"]; exists {
			ctx.Set("userRole", role)
		}

		ctx.Next()
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
