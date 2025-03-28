package routes

import (
	"github.com/ManoVikram/Threads-Knock-Off-API/handlers"
	"github.com/ManoVikram/Threads-Knock-Off-API/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// Get details of a specific user
	server.GET("/api/user/:id", handlers.GetUserDetailsHandler)

	// Protected routes
	protectedRoutes := server.Group("/api")
	protectedRoutes.Use(middlewares.AuthMiddleware())

	// Update the username of the user
	protectedRoutes.PATCH("/user/:id/username", handlers.UpdateUsernameHandler)
}
