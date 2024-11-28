package router

import (
	"movies/handler"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, UserHandler *handler.UsersHandler) {
	user := r.Group("/user")
	{
		user.POST("/register", UserHandler.Register) // Register a new user
		user.POST("/login", UserHandler.Login)       // Login an existing user
	}
}
