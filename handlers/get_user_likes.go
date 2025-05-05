package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/ManoVikram/Threads-Knock-Off-API/lib"
	"github.com/gin-gonic/gin"
)

func GetUserLikesHandler(c *gin.Context) {
	var userID string
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID = userIDValue.(string)

	// Query to get posts liked by the user, joined with post and user details
	rows, err := database.DB.Query(`
		SELECT 
			p.id, p.content, p.created_at, 
			u.id, u.name, u.username, u.image,
			p.likes_count, p.retweets_count, p.comments_count
		FROM likes l
		JOIN posts p ON l.post_id = p.id
		JOIN users u ON p.user_id = u.id
		WHERE l.user_id = $1
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch liked posts"})
		return
	}
	defer rows.Close()

	var likedPosts []gin.H

	for rows.Next() {
		var postID, content, createdAt string
		var userID, name, username string
		var image sql.NullString
		var likesCount, commentsCount, retweetsCount int

		err := rows.Scan(
			&postID, &content, &createdAt,
			&userID, &name, &username,
			&image,
			&likesCount, &commentsCount, &retweetsCount,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}

		likedPosts = append(likedPosts, gin.H{
			"id":             postID,
			"content":        content,
			"created_at":     createdAt,
			"likes_count":    likesCount,
			"comments_count": commentsCount,
			"retweet_count":  retweetsCount,
			"user": gin.H{
				"id":       userID,
				"name":     name,
				"username": username,
				"image":    lib.StringOrEmpty(image),
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"liked_posts": likedPosts,
	})
}
