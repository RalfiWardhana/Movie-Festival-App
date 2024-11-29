package router

import (
	"movies/handler"
	"movies/middleware"

	"github.com/gin-gonic/gin"
)

func StatsRoutes(r *gin.RouterGroup, StatsHandler *handler.StatsHandler) {
	stats := r.Group("/stats")
	{
		stats.POST("/:movie_id/view", middleware.AuthMiddleware(), StatsHandler.TrackView)                                                      // Authenticated users can track views
		stats.POST("/:movie_id/vote", middleware.AuthMiddleware(), StatsHandler.VoteMovie)                                                      // Authenticated users can vote
		stats.POST("/:movie_id/unvote", middleware.AuthMiddleware(), StatsHandler.UnvoteMovie)                                                  // Authenticated users can unvote
		stats.GET("/most-viewed-genre-movie", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), StatsHandler.GetMostViewedStats) // Admin can view stats
		stats.GET("/most-voted-genre-movie", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), StatsHandler.GetMostVotedStats)
		stats.POST("/:movie_id/trace", middleware.AuthMiddleware(), StatsHandler.TraceViewership)
		stats.GET("/user/voted-movies", middleware.AuthMiddleware(), StatsHandler.GetUserVotedMovies) // Authenticated users can view their votes
	}
}
