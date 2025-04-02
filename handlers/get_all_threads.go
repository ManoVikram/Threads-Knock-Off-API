package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/ManoVikram/Threads-Knock-Off-API/models"
	"github.com/gin-gonic/gin"
)

func GetAllThreadsHandler(c *gin.Context) {
	query := `
		SELECT 
			p.id, p.user_id, p.content, p.parent_id, 
			p.likes_count, p.retweets_count, p.comments_count, p.created_at,
			u.id, u.name, u.email, u."emailVerified", u.image, u.username, u.bio
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC;
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}
	defer rows.Close()

	var posts []gin.H

	for rows.Next() {
		var post models.Post
		var user models.User

		err := rows.Scan(
			&post.ID, &post.UserID, &post.Content, &post.ParentID,
			&post.LikesCount, &post.RetweetsCount, &post.CommentsCount, &post.CreatedAt,
			&user.ID, &user.Name, &user.Email, &user.EmailVerified, &user.Image, &user.Username, &user.Bio,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning posts"})
			return
		}

		posts = append(posts, gin.H{
			"id":             post.ID,
			"user_id":        post.UserID,
			"content":        post.Content,
			"parent_id":      post.ParentID,
			"likes_count":    post.LikesCount,
			"retweets_count": post.RetweetsCount,
			"comments_count": post.CommentsCount,
			"created_at":     post.CreatedAt,
			"user": gin.H{
				"id":             user.ID,
				"name":           user.Name,
				"email":          user.Email,
				"email_verified": user.EmailVerified.Time.String(),
				"image":          stringOrEmpty(user.Image),
				"username":       stringOrEmpty(user.Username),
				"bio":            stringOrEmpty(user.Bio),
			},
		})
	}

	c.JSON(http.StatusOK, posts)
}

// Helper function to return an empty string if sql.NullString is invalid
func stringOrEmpty(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
