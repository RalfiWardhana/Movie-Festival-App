package handler

import (
	"fmt"
	"log"
	"movies/middleware"
	"movies/model"
	"movies/usecase"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type MoviesHandler struct {
	MoviesUsecase *usecase.MoviesUseCase
}

func NewMoviesHandler(MoviesUseCase *usecase.MoviesUseCase) *MoviesHandler {
	return &MoviesHandler{MoviesUsecase: MoviesUseCase}
}

// Create handles movie creation with form data and file upload.
func (h *MoviesHandler) Create(c *gin.Context) {
	// Parse form data from the request.
	title := c.PostForm("title")
	description := c.PostForm("description")
	duration := c.PostForm("duration")
	artist := c.PostForm("artist")
	genreIDStr := c.PostForm("genre_id")

	// Initialize a slice to track missing fields for validation.
	var missingFields []string

	// Check for missing fields and append their names to the slice.
	if title == "" {
		missingFields = append(missingFields, "title")
	}
	if description == "" {
		missingFields = append(missingFields, "description")
	}
	if duration == "" {
		missingFields = append(missingFields, "duration")
	}
	if artist == "" {
		missingFields = append(missingFields, "artist")
	}
	if genreIDStr == "" {
		missingFields = append(missingFields, "genre_id")
	}

	// Return an error response if there are missing fields.
	if len(missingFields) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields", "fields": missingFields})
		return
	}

	// Convert the genre_id from a string to an integer.
	var genreID *int
	if genreIDStr != "" {
		gid, err := strconv.Atoi(genreIDStr)
		if err == nil {
			genreID = &gid
		}
	}

	// Retrieve the uploaded file from the form data.
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Validate the file size, allowing a maximum of 50 MB.
	const maxFileSize = 50 << 20 // 50 MB.
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the 50 MB limit"})
		return
	}

	// Define a list of allowed MIME types for video files.
	allowedMimeTypes := []string{
		"video/mp4",
		"video/x-matroska", // Supports .mkv format.
		"video/x-msvideo",  // Supports .avi format.
		"video/webm",
	}

	// Open and read the file header to validate its MIME type.
	fileHeader, _ := file.Open()
	defer fileHeader.Close()

	buffer := make([]byte, 512)
	_, err = fileHeader.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Detect the file type and ensure it matches allowed MIME types.
	fileType := http.DetectContentType(buffer)
	validFile := false
	for _, mimeType := range allowedMimeTypes {
		if fileType == mimeType {
			validFile = true
			break
		}
	}

	// Return an error response if the file type is not valid.
	if !validFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only movie/video files are allowed"})
		return
	}

	// Ensure the "uploads" directory exists, creating it if necessary.
	uploadPath := "uploads"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.Mkdir(uploadPath, os.ModePerm)
	}

	// Generate a unique filename for the uploaded file.
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	filePath := filepath.Join(uploadPath, filename)

	// Save the uploaded file to the designated directory.
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Create a movie object with the provided form data and generated file path.
	movie := &model.Movies{
		Title:       title,
		Description: description,
		Duration:    duration,
		Artist:      artist,
		GenreID:     genreID,
		WatchURL:    fmt.Sprintf("/%s", filePath),
	}

	// Use the use case layer to save the movie to the database.
	movieID, err := h.MoviesUsecase.Create(movie)
	if err != nil {
		if strings.Contains(err.Error(), "genre not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track viewership"})
		}
		return
	}

	// Return a success response with the new movie's details.
	c.JSON(http.StatusCreated, gin.H{
		"message":   "Movie uploaded successfully",
		"movie_id":  movieID,
		"watch_url": movie.WatchURL,
	})
}

// UpdateMovie handles the update of a movie record based on provided form data.
func (h *MoviesHandler) UpdateMovie(c *gin.Context) {
	// Parse movie ID from URL and handle invalid ID error
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	// Prepare a map to hold fields to be updated
	updates := make(map[string]interface{})
	var genreId int

	// Collect form data for fields like title, description, duration, artist, and genre
	// Only add non-empty fields to the updates map
	if title := c.PostForm("title"); title != "" {
		updates["title"] = title
	}
	if description := c.PostForm("description"); description != "" {
		updates["description"] = description
	}
	if durationStr := c.PostForm("duration"); durationStr != "" {
		if duration, err := strconv.Atoi(durationStr); err == nil {
			updates["duration"] = duration
		}
	}
	if artist := c.PostForm("artist"); artist != "" {
		updates["artist"] = artist
	}
	if genreIDStr := c.PostForm("genre_id"); genreIDStr != "" {
		if genreID, err := strconv.Atoi(genreIDStr); err == nil {
			updates["genre_id"] = genreID
			genreId = genreID
		}

	}

	// Handle file upload: validate size, type, and save the file
	// Update the movie's watch_url with the saved file path
	if file, err := c.FormFile("file"); err == nil {
		// Validate file size (max 50 MB)
		const maxFileSize = 50 << 20 // 50 MB
		if file.Size > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the 50 MB limit"})
			return
		}

		// Validate file type (only movie/video files allowed)
		allowedMimeTypes := []string{
			"video/mp4", "video/x-matroska", "video/x-msvideo", "video/webm",
		}

		fileHeader, _ := file.Open()
		defer fileHeader.Close()

		// Check the file type
		buffer := make([]byte, 512)
		_, err = fileHeader.Read(buffer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}

		fileType := http.DetectContentType(buffer)
		validFile := false
		for _, mimeType := range allowedMimeTypes {
			if fileType == mimeType {
				validFile = true
				break
			}
		}

		if !validFile {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only movie/video files are allowed"})
			return
		}

		// Ensure upload directory exists and save the file
		uploadPath := "uploads"
		if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
			os.Mkdir(uploadPath, os.ModePerm)
		}

		// Generate a unique filename and save the uploaded file
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		filePath := filepath.Join(uploadPath, filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Add the watch_url to the update fields
		updates["watch_url"] = fmt.Sprintf("/%s", filePath)
	}

	// If no fields were provided for update, return an error
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	// Perform the update using the MoviesUsecase and handle any errors
	err = h.MoviesUsecase.Update(movieID, genreId, updates)
	if err != nil {
		if strings.Contains(err.Error(), "genre not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
		} else if strings.Contains(err.Error(), "movie not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track viewership"})
		}
		return
	}

	// Return success message after a successful update
	c.JSON(http.StatusOK, gin.H{
		"message": "Movie updated successfully",
	})
}

// GetAllMoviesWithPagination handles fetching a list of movies with pagination
func (h *MoviesHandler) GetAllMoviesWithPagination(c *gin.Context) {
	// Retrieve the 'page' and 'limit' parameters from the query string (URL query parameters).
	// The DefaultQuery function is used to set default values for these parameters in case they are not provided by the client.
	page, err := strconv.Atoi(c.DefaultQuery("page", "1")) // Default page is 1.
	if err != nil {
		// If the page parameter is not a valid integer, return a 400 Bad Request response.
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default limit is 10.
	if err != nil {
		// If the limit parameter is not a valid integer, return a 400 Bad Request response.
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	// Call the MoviesUsecase to fetch the list of movies with pagination.
	// The use case function `GetAllMoviesWithPagination` takes the page and limit as arguments.
	movies, err := h.MoviesUsecase.GetAllMoviesWithPagination(page, limit)
	if err != nil {
		// If an error occurs while fetching the movies, return a 500 Internal Server Error response with the error message.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send the fetched movies data as a JSON response with a 200 OK status.
	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"messages": "Success to get Movies",
		"data":     len(movies),
		"movies":   movies,
	})
}

// SearchMovies handles the request to search for movies based on query parameters
func (h *MoviesHandler) SearchMovies(c *gin.Context) {
	// Get query parameters artist and genre_id
	title := c.Query("title")
	description := c.Query("description")
	artist := c.Query("artist")
	genreIDStr := c.DefaultQuery("genre_id", "0")
	genreID, err := strconv.Atoi(genreIDStr)
	if err != nil {
		genreID = 0
	}

	// Call usecase to search for movies
	movies, err := h.MoviesUsecase.SearchMovies(title, description, artist, genreID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return search results
	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"messages": "Success to get Movies",
		"movies":   movies,
	})
}

// TrackView tracks the view of a movie, with userID taken from the JWT token
func (h *MoviesHandler) TrackView(c *gin.Context) {
	// Extract the movie_id from the URL parameters
	movieIDStr := c.Param("movie_id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

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

	// Create the MovieView record with the current timestamp
	movieView := &model.MovieView{
		MovieID:  movieID,
		UserID:   userID,
		ViewedAt: time.Now().Format(time.RFC3339),
	}

	log.Println("Movie View: ", movieView)

	// Call the usecase to check and save view data
	hasViewed, err := h.MoviesUsecase.TrackMovieView(movieView)
	if err != nil {
		if strings.Contains(err.Error(), "movie not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track viewership"})
		}
		return
	}

	// Respond with a message if the user has already viewed the movie
	if hasViewed {
		c.JSON(http.StatusOK, gin.H{"message": "User has already viewed this movie"})
		return
	}

	// Respond with a success message for new views
	c.JSON(http.StatusOK, gin.H{"message": "Viewership tracked successfully"})
}
