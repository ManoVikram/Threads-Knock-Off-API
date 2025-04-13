package handlers

import (
	"net/http"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/ManoVikram/Threads-Knock-Off-API/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PostThreadHandler(c *gin.Context) {
	var post models.Post

	// Bind the incoming JSON to the Post struct
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Ensure content is not empty
	if post.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content cannot be empty"})
		return
	}

	// Get user ID from context (auth middleware should set this)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert string to UUID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Generate a new UUID for the post
	post.UserID = userID

	// Check if this is a reply (has a valid parent_id)
	if post.ParentID != nil {
		query := `
			INSERT INTO posts (user_id, content, parent_id)
			VALUES ($1, $2, $3)
		`
		_, err = database.DB.Exec(query, post.UserID, post.Content, post.ParentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reply"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Reply posted successfully"})
	} else {
		// Regular thread post
		query := `
			INSERT INTO posts (user_id, content)
			VALUES ($1, $2)
		`
		_, err = database.DB.Exec(query, post.UserID, post.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})
	}
}
