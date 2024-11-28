package router

import (
	"movies/handler"

	"github.com/gin-gonic/gin"
)

func Router(MoviesHandler *handler.MoviesHandler, StatsHandler *handler.StatsHandler, UserHandler *handler.UsersHandler) *gin.Engine {
	r := gin.Default()

	// Group routes
	api := r.Group("/api/v1")
	{
		UserRoutes(api, UserHandler)
		MovieRoutes(api, MoviesHandler)
		StatsRoutes(api, StatsHandler)
	}

	return r
}
