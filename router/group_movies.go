package router

import (
	"movies/handler"
	"movies/middleware"

	"github.com/gin-gonic/gin"
)

func MovieRoutes(r *gin.RouterGroup, MoviesHandler *handler.MoviesHandler) {
	movies := r.Group("/movies")
	{
		movies.POST("", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), MoviesHandler.Create)               // Admin can create movies
		movies.PUT("/:movie_id", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), MoviesHandler.UpdateMovie) // Admin can update movies
		movies.GET("", MoviesHandler.GetAllMoviesWithPagination)                                                             // Public route to get all movies
		movies.GET("/search", MoviesHandler.SearchMovies)                                                                    // Public route to search movies
	}
}
