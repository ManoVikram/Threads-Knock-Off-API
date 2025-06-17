package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/ManoVikram/Threads-Knock-Off-API/models"
	"github.com/gin-gonic/gin"
)

func GetUserDetailsHandler(c *gin.Context) {
	userID := c.Query("userid")
	username := c.Query("username")

	if userID == "" && username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either id or username must be provided"})
		return
	}

	var user models.User
	var query string
	var args []any

	if userID != "" && username != "" {
		query = `SELECT id, name, email, "emailVerified", image, username, follower_count, following_count, bio FROM users WHERE id = $1 AND username = $2`
		args = append(args, userID, username)
	} else if userID != "" {
		query = `SELECT id, name, email, "emailVerified", image, username, follower_count, following_count, bio FROM users WHERE id = $1`
		args = append(args, userID)
	} else {
		query = `SELECT id, name, email, "emailVerified", image, username, follower_count, following_count, bio FROM users WHERE username = $1`
		args = append(args, username)
	}

	err := database.DB.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.EmailVerified,
		&user.Image,
		&user.Username,
		&user.FollowerCount,
		&user.FollowingCount,
		&user.Bio,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Convert sql.Null types to safe values
	response := gin.H{
		"id":              user.ID,
		"name":            user.Name,
		"email":           user.Email,
		"follower_count":  user.FollowerCount,
		"following_count": user.FollowingCount,
	}

	if user.EmailVerified.Valid {
		response["emailVerified"] = user.EmailVerified.Time.Format(time.RFC3339)
	} else {
		response["emailVerified"] = ""
	}

	response["image"] = user.Image.String
	response["username"] = user.Username.String
	response["bio"] = user.Bio.String

	c.JSON(http.StatusOK, response)
}
