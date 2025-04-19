package handlers

import (
	"net/http"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/ManoVikram/Threads-Knock-Off-API/lib"
	"github.com/ManoVikram/Threads-Knock-Off-API/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetThreadHandler(c *gin.Context) {
	postIDParam := c.Param("id")
	postID, err := uuid.Parse(postIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var userID string
	userIDValue, exists := c.Get("userID")
	if exists {
		userID = userIDValue.(string)
	}

	// Fetch post and user details
	query := `
		SELECT 
			p.id, p.user_id, p.content, p.parent_id, 
			p.likes_count, p.retweets_count, p.comments_count, p.created_at,
			u.id, u.name, u.email, u."emailVerified", u.image, u.username, u.bio
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $1
	`

	var post models.Post
	var author models.User

	err = database.DB.QueryRow(query, postID).Scan(
		&post.ID, &post.UserID, &post.Content, &post.ParentID,
		&post.LikesCount, &post.RetweetsCount, &post.CommentsCount, &post.CreatedAt,
		&author.ID, &author.Name, &author.Email, &author.EmailVerified, &author.Image, &author.Username, &author.Bio,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if the user has liked the post
	var hasLiked bool
	if exists {
		likeCheck := `SELECT EXISTS (SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2)`
		err = database.DB.QueryRow(likeCheck, userID, post.ID).Scan(&hasLiked)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check like status"})
			return
		}
	}

	// Fetch comments on this post
	commentQuery := `
		SELECT 
			p.id, p.user_id, p.content, p.likes_count, p.retweets_count, p.comments_count, p.created_at,
			u.id, u.name, u.email, u."emailVerified", u.image, u.username, u.bio
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.parent_id = $1
		ORDER BY p.created_at ASC
	`

	rows, err := database.DB.Query(commentQuery, post.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}
	defer rows.Close()

	var comments []gin.H

	for rows.Next() {
		var comment models.Post
		var commenter models.User
		err := rows.Scan(
			&comment.ID, &comment.UserID, &comment.Content, &comment.LikesCount, &comment.RetweetsCount, &comment.CommentsCount, &comment.CreatedAt,
			&commenter.ID, &commenter.Name, &commenter.Email, &commenter.EmailVerified, &commenter.Image, &commenter.Username, &commenter.Bio,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning comment"})
			return
		}

		// Check if user liked the comment
		var commentLiked bool
		if exists {
			likeCheck := `SELECT EXISTS (SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2)`
			err = database.DB.QueryRow(likeCheck, userID, comment.ID).Scan(&commentLiked)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check like on comment"})
				return
			}
		}

		comments = append(comments, gin.H{
			"id":             comment.ID,
			"user_id":        comment.UserID,
			"content":        comment.Content,
			"likes_count":    comment.LikesCount,
			"retweets_count": comment.RetweetsCount,
			"comments_count": comment.CommentsCount,
			"created_at":     comment.CreatedAt,
			"liked_by_user":  commentLiked,
			"user": gin.H{
				"id":             commenter.ID,
				"name":           commenter.Name,
				"email":          commenter.Email,
				"email_verified": commenter.EmailVerified.Time.String(),
				"image":          lib.StringOrEmpty(commenter.Image),
				"username":       lib.StringOrEmpty(commenter.Username),
				"bio":            lib.StringOrEmpty(commenter.Bio),
			},
		})
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
			"id":             author.ID,
			"name":           author.Name,
			"email":          author.Email,
			"email_verified": author.EmailVerified.Time.String(),
			"image":          lib.StringOrEmpty(author.Image),
			"username":       lib.StringOrEmpty(author.Username),
			"bio":            lib.StringOrEmpty(author.Bio),
		},
		"comments": comments,
	})
}
