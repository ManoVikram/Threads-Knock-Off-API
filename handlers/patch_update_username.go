package handlers

import (
	"net/http"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/gin-gonic/gin"
)

// UpdateUsernameHandler updates the username for a specific user
func UpdateUsernameHandler(c *gin.Context) {
	// Extract userID from middleware context
	authenticatedUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
		return
	}

	// Extract user ID from URL param
	userID := c.Param("id")

	// Ensure the authenticated user is the same as the one being updated
	if authenticatedUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own username"})
		return
	}

	// Parse request body
	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Ensure username is provided
	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username cannot be empty"})
		return
	}

	// Check if the username already exists
	err := database.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, req.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	// Update the username
	_, err = database.DB.Exec(`UPDATE users SET username = $1 WHERE id = $2`, req.Username, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{"message": "Username updated successfully"})
}
