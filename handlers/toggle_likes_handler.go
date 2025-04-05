package handlers

import (
	"net/http"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ToggleLikePostHandler(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	postID := c.Param("id")

	// Check if like already exists
	var existsFlag bool
	err = database.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2
		)
	`, userID, postID).Scan(&existsFlag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if existsFlag {
		// Unlike
		_, err = database.DB.Exec(`
			DELETE FROM likes WHERE user_id = $1 AND post_id = $2
		`, userID, postID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike post"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Post unliked"})
		return
	} else {
		// Like
		_, err = database.DB.Exec(`
			INSERT INTO likes (user_id, post_id) VALUES ($1, $2)
		`, userID, postID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Post liked"})
		return
	}
}
