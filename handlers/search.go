package handlers

import (
	"net/http"
	"strings"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/ManoVikram/Threads-Knock-Off-API/models"
	"github.com/gin-gonic/gin"
)

func SearchHandler(c *gin.Context) {
	query := c.Query("q") // Get search text from query param e.g., /search?q=hello

	if strings.TrimSpace(query) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query cannot be empty"})
		return
	}

	// Search Posts
	postRows, err := database.DB.Query(`
		SELECT id, user_id, content, created_at
		FROM posts
		WHERE to_tsvector('english', content) @@ plainto_tsquery('english', $1)
	`, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search posts"})
		return
	}
	defer postRows.Close()

	var posts []models.Post
	for postRows.Next() {
		var post models.Post
		err := postRows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read post data"})
			return
		}
		posts = append(posts, post)
	}

	// Search Users
	userRows, err := database.DB.Query(`
		SELECT id, name, email, username, bio
		FROM users
		WHERE to_tsvector(
			'english',
			COALESCE(name, '') || ' ' || COALESCE(username, '')
		) @@ plainto_tsquery('english', $1)
	`, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}
	defer userRows.Close()

	var users []models.User
	for userRows.Next() {
		var user models.User
		err := userRows.Scan(&user.ID, &user.Name, &user.Email, &user.Username, &user.Bio)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user data"})
			return
		}
		users = append(users, user)
	}

	// Send both users and posts as response
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"users": users,
	})
}
