package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/ManoVikram/Threads-Knock-Off-API/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetThreadHandler(c *gin.Context) {
	threadIDParam := c.Param("id")
	threadID, err := uuid.Parse(threadIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid thread ID"})
		return
	}

	query := `
		SELECT 
			p.id, p.user_id, p.content, p.parent_id,
			p.likes_count, p.retweets_count, p.comments_count, p.created_at,
			u.id, u.name, u.email, u."emailVerified", u.image, u.username, u.bio
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $1;
	`

	row := database.DB.QueryRow(query, threadID)

	var post models.Post
	var user models.User

	err = row.Scan(
		&post.ID, &post.UserID, &post.Content, &post.ParentID,
		&post.LikesCount, &post.RetweetsCount, &post.CommentsCount, &post.CreatedAt,
		&user.ID, &user.Name, &user.Email, &user.EmailVerified, &user.Image, &user.Username, &user.Bio,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Thread not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve thread"})
		}
		return
	}

	// Optional: Check if user is logged in and has liked this post
	var hasLiked bool = false
	if userIDRaw, exists := c.Get("userID"); exists {
		if userID, err := uuid.Parse(userIDRaw.(string)); err == nil {
			likeQuery := `
				SELECT EXISTS (
					SELECT 1 FROM likes WHERE post_id = $1 AND user_id = $2
				)
			`
			err := database.DB.QueryRow(likeQuery, post.ID, userID).Scan(&hasLiked)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking like status"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             post.ID,
		"user_id":        post.UserID,
		"content":        post.Content,
		"parent_id":      post.ParentID,
		"likes_count":    post.LikesCount,
		"retweets_count": post.RetweetsCount,
		"comments_count": post.CommentsCount,
		"created_at":     post.CreatedAt,
		"liked_by_user":  hasLiked,
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
