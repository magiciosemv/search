package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware returns a gin middleware that validates Bearer token.
// If apiKey is empty, auth is disabled (dev mode).
func AuthMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiKey == "" {
			c.Next()
			return
		}

		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
				"code":  "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format, use: Bearer <key>",
				"code":  "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		if parts[1] != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid api key",
				"code":  "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
