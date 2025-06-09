package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/ManoVikram/Threads-Knock-Off-API/lib"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserPostsHandler(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	// Get the user's ID from the username
	var userID string
	row := database.DB.QueryRow("SELECT id FROM users WHERE username = $1", username)
	if err := row.Scan(&userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get the currently logged-in user's ID
	currentUserIDStr, userLoggedIn := c.Get("userID")

	// Pre-fetch all post IDs liked by the current user
	likedPostIDs := make(map[string]bool)
	if userLoggedIn {
		currentUserID, err := uuid.Parse(currentUserIDStr.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}

		likeRows, err := database.DB.Query(`SELECT post_id FROM likes WHERE user_id = $1`, currentUserID)
		if err == nil {
			defer likeRows.Close()
			for likeRows.Next() {
				var likedPostID string
				if err := likeRows.Scan(&likedPostID); err == nil {
					likedPostIDs[likedPostID] = true
				}
			}
		}
	}

	// Query to fetch posts and replies
	rows, err := database.DB.Query(`
		SELECT 
			p.id, p.content, p.created_at, p.parent_id,
			u.id, u.name, u.username, u.bio, u.image,
			p.likes_count, p.retweets_count, p.comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = $1
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var originalPosts []gin.H
	var replyPosts []gin.H

	for rows.Next() {
		var postID, content, createdAt string
		var parentID sql.NullString
		var uID, name, username string
		var bio, image sql.NullString
		var likesCount, commentsCount, retweetCount int

		err := rows.Scan(
			&postID, &content, &createdAt, &parentID,
			&uID, &name, &username, &bio, &image,
			&likesCount, &retweetCount, &commentsCount,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read post data"})
			return
		}

		postData := gin.H{
			"id":             postID,
			"content":        content,
			"created_at":     createdAt,
			"parent_id":      lib.StringOrEmpty(parentID),
			"likes_count":    likesCount,
			"comments_count": commentsCount,
			"retweets_count": retweetCount,
			"liked_by_user":  likedPostIDs[postID],
			"user": gin.H{
				"id":       uID,
				"name":     name,
				"username": username,
				"bio":      lib.StringOrEmpty(bio),
				"image":    lib.StringOrEmpty(image),
			},
		}

		if parentID.Valid {
			replyPosts = append(replyPosts, postData)
		} else {
			originalPosts = append(originalPosts, postData)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"posts":   originalPosts,
		"replies": replyPosts,
	})
}
