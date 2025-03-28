package routes

import (
	"github.com/ManoVikram/Threads-Knock-Off-API/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// Get details of a specific user
	server.GET("api/user/:id", handlers.GetUserDetailsHandler)
}
