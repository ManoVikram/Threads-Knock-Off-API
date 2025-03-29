package middlewares

import (
	"net/http"
	"strings"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks if the request has a valid Bearer token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Extract the token (format: "Bearer <token>")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		sessionToken := parts[1]
		var userID string

		// Validate token against the database
		query := `SELECT "userId" FROM sessions WHERE "sessionToken" = $1`
		err := database.DB.QueryRow(query, sessionToken).Scan(&userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store user ID in the context
		c.Set("userID", userID)
		c.Next()
	}
}
