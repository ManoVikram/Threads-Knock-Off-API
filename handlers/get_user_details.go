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
	userID := c.Param("id") // Assuming ID is passed in URL

	var user models.User
	query := `SELECT id, name, email, "emailVerified", image, username, bio FROM users WHERE id = $1`
	err := database.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.EmailVerified,
		&user.Image,
		&user.Username,
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

	// Convert sql.NullString and sql.NullTime to regular JSON-friendly types
	response := gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}

	if user.EmailVerified.Valid {
		emailVerified := user.EmailVerified.Time.Format(time.RFC3339)
		response["emailVerified"] = emailVerified
	} else {
		response["emailVerified"] = ""
	}
	response["image"] = user.Image.String
	response["username"] = user.Username.String
	response["bio"] = user.Bio.String

	c.JSON(http.StatusOK, response)
}
