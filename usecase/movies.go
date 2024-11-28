package usecase

import (
	"errors"
	"fmt"
	"movies/model"
	"movies/repository"
)

type MoviesUseCase struct {
	MoviesRepo *repository.MoviesRepo
}

func NewMoviesUseCase(MoviesRepo *repository.MoviesRepo) *MoviesUseCase {
	return &MoviesUseCase{MoviesRepo: MoviesRepo}
}

// Create creates a new movie by calling the repository method.
func (uc *MoviesUseCase) Create(Movie *model.Movies) (*model.Movies, error) {
	// Check movies
	movieExists, err := uc.MoviesRepo.GenreExists(*Movie.GenreID) // Checks if genre exists
	if err != nil {
		return nil, fmt.Errorf("failed to validate genre existence: %w", err)
	}
	if !movieExists {
		return nil, fmt.Errorf("genre not found") // Returns error if genre doesn't exist
	}
	return uc.MoviesRepo.Create(Movie) // Calls repository to create movie
}

// Update updates an existing movie using the repository method.
func (uc *MoviesUseCase) Update(movieID int, genreID int, updates map[string]interface{}) error {
	// Check genres
	genreExists, err := uc.MoviesRepo.GenreExists(genreID) // Checks if genre exists
	if err != nil {
		return fmt.Errorf("failed to validate genre existence: %w", err)
	}
	if !genreExists {
		return fmt.Errorf("genre not found") // Returns error if genre doesn't exist
	}

	// Check movies
	movieExists, err := uc.MoviesRepo.MovieExists(movieID) // Checks if movie exists
	if err != nil {
		return fmt.Errorf("failed to validate movie existence: %w", err)
	}
	if !movieExists {
		return fmt.Errorf("movie not found") // Returns error if movie doesn't exist
	}

	return uc.MoviesRepo.Update(movieID, updates) // Calls repository to update movie
}

// GetAllMoviesWithPagination retrieves movies with pagination by calling the repository.
func (uc *MoviesUseCase) GetAllMoviesWithPagination(page, limit int) ([]model.Movies, error) {
	if page <= 0 || limit <= 0 {
		return nil, errors.New("invalid page or limit") // Validates page and limit
	}
	return uc.MoviesRepo.GetAllMoviesWithPagination(page, limit) // Calls repository to get movies with pagination
}

// SearchMovies calls the repository to search movies by artist, genre_id, or both
func (uc *MoviesUseCase) SearchMovies(title string, description string, artist string, genreID int) ([]model.Movies, error) {
	return uc.MoviesRepo.SearchMovies(title, description, artist, genreID) // Calls repository to search movies
}

// TrackMovieView saves the viewership record into the database
func (uc *MoviesUseCase) TrackMovieView(view *model.MovieView) (bool, error) {
	// Check movies
	movieExists, err := uc.MoviesRepo.MovieExists(view.MovieID) // Checks if movie exists
	if err != nil {
		return false, fmt.Errorf("failed to validate movie existence: %w", err)
	}
	if !movieExists {
		return false, fmt.Errorf("movie not found") // Returns error if movie doesn't exist
	}

	// Check if user has already viewed the movie
	hasViewed, err := uc.MoviesRepo.HasViewed(view.MovieID, view.UserID) // Checks if user has viewed the movie
	if err != nil {
		return false, err
	}

	// If user has already viewed, return true
	if hasViewed {
		return true, nil
	}

	// Save the new view record to database
	err = uc.MoviesRepo.SaveMovieView(view) // Calls repository to save view
	if err != nil {
		return false, err
	}

	return false, nil
}
