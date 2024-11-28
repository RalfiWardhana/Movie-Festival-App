package handler

import (
	"log"
	"movies/middleware"
	"movies/model"
	"movies/usecase"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	StatsUseCase *usecase.StatsUseCase
}

// Constructor untuk StatsHandler
func NewStatsHandler(statsUseCase *usecase.StatsUseCase) *StatsHandler {
	return &StatsHandler{
		StatsUseCase: statsUseCase,
	}
}

// GetMostViewedStats handles the request to retrieve the most viewed statistics.
func (h *StatsHandler) GetMostViewedStats(c *gin.Context) {
	// Call the use case to get the most viewed statistics data
	stats, err := h.StatsUseCase.GetMostViewedStats()

	// If there is an error in fetching statistics, return an internal server error response with the error message
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If no error, send a successful response (200 OK) with the statistics data as JSON
	c.JSON(http.StatusOK, stats)
}

// GetMostVotedStats handles the request to retrieve the most voted statistics.
func (h *StatsHandler) GetMostVotedStats(c *gin.Context) {
	// Call the use case to get the most voted statistics data
	stats, err := h.StatsUseCase.GetMostVotedStats()

	// If there is an error in fetching statistics, return an internal server error response with the error message
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If no error, send a successful response (200 OK) with the statistics data as JSON
	c.JSON(http.StatusOK, stats)
}

// VoteMovie handles voting for a movie
func (h *StatsHandler) VoteMovie(c *gin.Context) {
	// Extract the userID from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Assert claims to the expected type
	userClaims, ok := claims.(*middleware.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Retrieve the userID from the claims
	userID := userClaims.UserID

	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	response, err := h.StatsUseCase.VoteMovie(userID, movieID)
	if err != nil {
		if strings.Contains(err.Error(), "movie not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else if strings.Contains(err.Error(), "movie has not been viewed") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Movie has not been viewed"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track viewership"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// UnvoteMovie handles unvoting for a movie
func (h *StatsHandler) UnvoteMovie(c *gin.Context) {
	// Extract the userID from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Assert claims to the expected type
	userClaims, ok := claims.(*middleware.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Retrieve the userID from the claims
	userID := userClaims.UserID

	movieID, err := strconv.Atoi(c.Param("movie_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	response, err := h.StatsUseCase.UnvoteMovie(userID, movieID)
	if err != nil {
		if strings.Contains(err.Error(), "movie not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else if strings.Contains(err.Error(), "movie has not been viewed") {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Movie has not been viewed"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track viewership"})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUserVotedMovies handles the request to retrieve a list of user's voted movies.
func (h *StatsHandler) GetUserVotedMovies(c *gin.Context) {
	// Extract the userID from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Assert claims to the expected type
	userClaims, ok := claims.(*middleware.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Retrieve the userID from the claims
	userID := userClaims.UserID
	log.Println(userID)

	// Fetch the voted movies using the use case
	votedMovies, err := h.StatsUseCase.GetUserVotedMovies(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of voted movies
	c.JSON(http.StatusOK, gin.H{"voted_movies": votedMovies})
}

// TraceViewership handles tracking of viewership by watching duration
func (h *StatsHandler) TraceViewership(c *gin.Context) {
	// Extract the userID from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Assert claims to the expected type
	userClaims, ok := claims.(*middleware.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Retrieve the userID from the claims
	userID := userClaims.UserID

	// Parse movie_id and duration from the request
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var request model.RequestDuration
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	if request.Duration == 0 {
		c.JSON(400, gin.H{"error": "Bad request", "message": "field duration required"})
		return
	}

	// Call use case to update the viewership duration
	err = h.StatsUseCase.TraceViewership(userID, movieID, request.Duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Viewership duration tracked successfully"})
}
