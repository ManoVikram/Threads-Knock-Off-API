package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ToggleFollowHandler(c *gin.Context) {
	// 1. Get the logged-in user ID
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	followerID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// 2. Get the target username from URL param
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	// 3. Get the target user ID from the username
	var followingID uuid.UUID
	err = database.DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&followingID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	// 4. Prevent users from following themselves
	if followerID == followingID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}

	// 5. Check if already following
	var existsFlag bool
	err = database.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2
		)
	`, followerID, followingID).Scan(&existsFlag)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 6. Toggle follow/unfollow
	if existsFlag {
		// Unfollow
		_, err := database.DB.Exec(`
			DELETE FROM follows WHERE follower_id = $1 AND following_id = $2
		`, followerID, followingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"isFollowing": false})
	} else {
		// Follow
		_, err := database.DB.Exec(`
			INSERT INTO follows (follower_id, following_id) VALUES ($1, $2)
		`, followerID, followingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"isFollowing": true})
	}
}
