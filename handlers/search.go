package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/ManoVikram/Threads-Knock-Off-API/lib"
	"github.com/ManoVikram/Threads-Knock-Off-API/models"
	"github.com/gin-gonic/gin"
)

func SearchHandler(c *gin.Context) {
	query := c.Query("q")

	if strings.TrimSpace(query) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query cannot be empty"})
		return
	}

	// Search Posts and Join with User Info
	postRows, err := database.DB.Query(`
		SELECT p.id, p.user_id, p.content, p.created_at,
		       u.name, u.username, u.bio, u.image
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE to_tsvector('english', p.content) @@ plainto_tsquery('english', $1)
		ORDER BY p.created_at DESC
	`, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search posts"})
		return
	}
	defer postRows.Close()

	var posts []gin.H
	for postRows.Next() {
		var postID, userID, content string
		var createdAt string
		var name, username string
		var bio, image sql.NullString

		err := postRows.Scan(&postID, &userID, &content, &createdAt,
			&name, &username, &bio, &image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read post data"})
			return
		}

		posts = append(posts, gin.H{
			"id":         postID,
			"content":    content,
			"created_at": createdAt,
			"user": gin.H{
				"id":       userID,
				"name":     name,
				"username": username,
				"bio":      lib.StringOrEmpty(bio),
				"image":    lib.StringOrEmpty(image),
			},
		})
	}

	// Search Users
	userRows, err := database.DB.Query(`
		SELECT id, name, email, username, bio, image
		FROM users
		WHERE to_tsvector('english', COALESCE(name, '') || ' ' || COALESCE(username, '')) @@ plainto_tsquery('english', $1)
	`, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}
	defer userRows.Close()

	var users []gin.H
	for userRows.Next() {
		var user models.User
		var username, bio, image sql.NullString

		err := userRows.Scan(&user.ID, &user.Name, &user.Email, &username, &bio, &image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user data"})
			return
		}

		users = append(users, gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"username": lib.StringOrEmpty(username),
			"bio":      lib.StringOrEmpty(bio),
			"image":    lib.StringOrEmpty(image),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"users": users,
	})
}
