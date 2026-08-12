package middlewares

import (
	"github.com/gin-gonic/gin"
)

func SetRoleMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("targetRole", role)
		c.Next()
	}
}
