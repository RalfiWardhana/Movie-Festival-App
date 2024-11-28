package router

import (
	"movies/handler"
	"movies/middleware"

	"github.com/gin-gonic/gin"
)

func Router(MoviesHandler *handler.MoviesHandler, StatsHandler *handler.StatsHandler, UserHandler *handler.UsersHandler) *gin.Engine {
	r := gin.Default()

	// User Routes
	r.POST("/register", UserHandler.Register) // Register a new user
	r.POST("/login", UserHandler.Login)       // Login an existing user

	// Movie Routes
	r.POST("/movies", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), MoviesHandler.Create)               // Admin can create movies
	r.PUT("/movies/:movie_id", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), MoviesHandler.UpdateMovie) // Admin can update movies
	r.GET("/movies", MoviesHandler.GetAllMoviesWithPagination)                                                             // Public route to get all movies
	r.GET("/movies/search", MoviesHandler.SearchMovies)                                                                    // Public route to search movies

	// Stats Routes
	r.POST("/stats/:movie_id/view", middleware.AuthMiddleware(), MoviesHandler.TrackView)                                                     // Authenticated users can track views
	r.POST("/stats/:movie_id/vote", middleware.AuthMiddleware(), middleware.RoleMiddleware("user"), StatsHandler.VoteMovie)                   // Authenticated users can vote
	r.POST("/stats/:movie_id/unvote", middleware.AuthMiddleware(), middleware.RoleMiddleware("user"), StatsHandler.UnvoteMovie)               // Authenticated users can unvote
	r.GET("/stats/most-viewed-genre-movie", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), StatsHandler.GetMostViewedStats) // Admin can view stats
	r.GET("/stats/most-voted-genre-movie", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), StatsHandler.GetMostVotedStats)
	r.POST("/stats/:movie_id/trace", middleware.AuthMiddleware(), middleware.RoleMiddleware("user"), StatsHandler.TraceViewership)
	r.GET("/stats/user/voted-movies", middleware.AuthMiddleware(), StatsHandler.GetUserVotedMovies) // Authenticated users can view their votes
	middleware.AuthMiddleware()
	return r
}
