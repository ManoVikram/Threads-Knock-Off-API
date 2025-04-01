package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/ManoVikram/Threads-Knock-Off-API/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllThreadsHandler(c *gin.Context) {
	query := `SELECT id, user_id, content, parent_id, likes_count, retweets_count, comments_count, created_at FROM posts ORDER BY created_at DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		var parentID sql.NullString

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Content,
			&parentID,
			&post.LikesCount,
			&post.RetweetsCount,
			&post.CommentsCount,
			&post.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing posts"})
			return
		}

		// Check if parent_id is valid before assigning
		if parentID.Valid {
			parsedParentID, parseErr := uuid.Parse(parentID.String)
			if parseErr == nil {
				post.ParentID = &parsedParentID
			}
		} else {
			post.ParentID = nil // Explicitly setting to nil
		}

		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
}
