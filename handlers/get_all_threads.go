package handlers

import (
	"net/http"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/ManoVikram/Threads-Knock-Off-API/lib"
	"github.com/ManoVikram/Threads-Knock-Off-API/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllThreadsHandler(c *gin.Context) {
	query := `
		SELECT 
			p.id, p.user_id, p.content, p.parent_id,
			p.likes_count, p.retweets_count, p.comments_count, p.created_at,
			u.id, u.name, u.email, u."emailVerified", u.image, u.username, u.bio
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.parent_id IS NULL
		ORDER BY p.created_at DESC;
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}
	defer rows.Close()

	var posts []gin.H

	userIDStr, userLoggedIn := c.Get("userID")
	var likedPostIDs map[string]bool

	if userLoggedIn {
		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}

		likeRows, err := database.DB.Query("SELECT post_id FROM likes WHERE user_id = $1", userID)
		if err == nil {
			defer likeRows.Close()
			likedPostIDs = make(map[string]bool)
			for likeRows.Next() {
				var postID string
				if err := likeRows.Scan(&postID); err == nil {
					likedPostIDs[postID] = true
				}
			}
		}
	}

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

		// Check if this post was liked by the logged-in user
		liked := false
		if userLoggedIn && likedPostIDs != nil {
			liked = likedPostIDs[post.ID.String()]
		}

		posts = append(posts, gin.H{
			"id":             post.ID,
			"user_id":        post.UserID,
			"content":        post.Content,
			"parent_id":      post.ParentID,
			"liked_by_user":  liked,
			"likes_count":    post.LikesCount,
			"retweets_count": post.RetweetsCount,
			"comments_count": post.CommentsCount,
			"created_at":     post.CreatedAt,
			"user": gin.H{
				"id":             user.ID,
				"name":           user.Name,
				"email":          user.Email,
				"email_verified": user.EmailVerified.Time.String(),
				"image":          lib.StringOrEmpty(user.Image),
				"username":       lib.StringOrEmpty(user.Username),
				"bio":            lib.StringOrEmpty(user.Bio),
			},
		})
	}

	c.JSON(http.StatusOK, posts)
}
