package routes

import (
	"github.com/ManoVikram/Threads-Knock-Off-API/handlers"
	"github.com/ManoVikram/Threads-Knock-Off-API/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// Get details of a specific user
	server.GET("/api/user/:id", handlers.GetUserDetailsHandler)

	// Get all the threads
	server.GET("/api/posts", middlewares.AuthMiddlewareLite(), handlers.GetAllThreadsHandler)

	// Get a single thread
	server.GET("/api/post/:id", middlewares.AuthMiddlewareLite(), handlers.GetThreadHandler)

	// Protected routes
	protectedRoutes := server.Group("/api")
	protectedRoutes.Use(middlewares.AuthMiddleware())

	// Update the username of the user
	protectedRoutes.PATCH("/user/:id/username", handlers.UpdateUsernameHandler)

	// Create a new thread / post
	protectedRoutes.POST("/posts", handlers.PostThreadHandler)

	// Like / Un-like a post
	protectedRoutes.POST("/posts/:id/like", handlers.ToggleLikePostHandler)
}
