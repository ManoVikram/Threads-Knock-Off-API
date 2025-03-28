package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ManoVikram/Threads-Knock-Off-API/database"
	"github.com/ManoVikram/Threads-Knock-Off-API/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserDetailsHandler(c *gin.Context) {
	userID := c.Param("id")

	// Parse UUID
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Query user details
	var user models.User
	query := `SELECT id, name, email, "emailVerified", image, username, bio FROM users WHERE id = $1`
	err = database.DB.QueryRow(query, parsedID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.EmailVerified,
		&user.Image,
		&user.Username,
		&user.Bio,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
