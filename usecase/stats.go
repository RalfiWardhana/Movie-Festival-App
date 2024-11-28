package usecase

import (
	"errors"
	"fmt"
	"movies/model"
	"movies/repository"
)

type StatsUseCase struct {
	StatsRepo *repository.StatsRepo
}

func NewStatsUseCase(StatsRepo *repository.StatsRepo) *StatsUseCase {
	return &StatsUseCase{StatsRepo: StatsRepo}
}

func (uc *StatsUseCase) GetMostViewedStats() (*model.StatsModelView, error) {
	// Fetch most viewed movies
	movies, err := uc.StatsRepo.GetMostViewedMovies()
	if err != nil {
		return nil, errors.New("failed to fetch most viewed movies")
	}

	// Fetch most viewed genres
	genres, err := uc.StatsRepo.GetMostViewedGenres()
	if err != nil {
		return nil, errors.New("failed to fetch most viewed genres")
	}

	// Return the combined statistics
	stats := &model.StatsModelView{
		MostViewedMovie: movies,
		MostViewedGenre: genres,
	}

	return stats, nil
}

func (uc *StatsUseCase) GetMostVotedStats() (*model.StatsModelVote, error) {
	// Fetch most voted movies
	movies, err := uc.StatsRepo.GetMostVotedMovies()
	if err != nil {
		return nil, errors.New("failed to fetch most voted movies")
	}

	// Fetch most viewed genres
	genres, err := uc.StatsRepo.GetMostViewedGenres()
	if err != nil {
		return nil, errors.New("failed to fetch most viewed genres")
	}

	// Return the combined statistics
	stats := &model.StatsModelVote{
		MostVotedMovie:  movies,
		MostViewedGenre: genres,
	}

	return stats, nil
}

// VoteMovie handles the logic for voting a movie
func (uc *StatsUseCase) VoteMovie(userID, movieID int) (map[string]interface{}, error) {
	// Check if movie exists
	movieExists, err := uc.StatsRepo.MovieExists(movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate movie existence: %w", err)
	}
	if !movieExists {
		return nil, fmt.Errorf("movie not found")
	}

	// Check if the movie has been viewed
	movieView, err := uc.StatsRepo.MovieViewExists(userID, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate movie view existence: %w", err)
	}
	if !movieView {
		return nil, fmt.Errorf("movie has not been viewed")
	}

	// Check if the user already voted for the movie
	voteExists, err := uc.StatsRepo.CheckVote(userID, movieID)
	if err != nil {
		return nil, err
	}

	// If the user already voted, return message
	if voteExists {
		return map[string]interface{}{
			"message": "User already voted for this movie",
			"status":  "unchanged",
		}, nil
	}

	// Add vote to the database
	err = uc.StatsRepo.AddVote(userID, movieID)
	if err != nil {
		return nil, err
	}

	// Return success message
	return map[string]interface{}{
		"message": "Movie voted successfully",
		"status":  "success",
	}, nil
}

// UnvoteMovie handles the logic for unvoting a movie
func (uc *StatsUseCase) UnvoteMovie(userID, movieID int) (map[string]interface{}, error) {
	// Check if movie exists
	movieExists, err := uc.StatsRepo.MovieExists(movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate movie existence: %w", err)
	}
	if !movieExists {
		return nil, fmt.Errorf("movie not found")
	}

	// Check if the movie has been viewed
	movieView, err := uc.StatsRepo.MovieViewExists(userID, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate movie view existence: %w", err)
	}
	if !movieView {
		return nil, fmt.Errorf("movie has not been viewed")
	}

	// Check if the user already unvoted the movie
	voteExists, err := uc.StatsRepo.CheckUnVote(userID, movieID)
	if err != nil {
		return nil, err
	}

	// If the user already unvoted, return message
	if voteExists {
		return map[string]interface{}{
			"message": "User already unvoted for this movie",
			"status":  "unchanged",
		}, nil
	}

	// Add unvote from the database
	err = uc.StatsRepo.AddUnVote(userID, movieID)
	if err != nil {
		return nil, err
	}

	// Return success message
	return map[string]interface{}{
		"message": "Movie unvote successfully",
		"status":  "success",
	}, nil
}

// GetUserVotedMovies retrieves all movies voted by a user.
func (uc *StatsUseCase) GetUserVotedMovies(userID int) ([]model.UserVotedMovie, error) {
	// Fetch voted movies for user
	votedMoviesData, err := uc.StatsRepo.GetUserVotedMovies(userID)
	if err != nil {
		return nil, errors.New("failed to fetch user's voted movies")
	}

	// Map the fetched data to UserVotedMovie struct
	var votedMovies []model.UserVotedMovie
	for _, movie := range votedMoviesData {
		votedMovies = append(votedMovies, model.UserVotedMovie{
			MovieID:  movie["movie_id"].(int),
			Title:    movie["title"].(string),
			IsLike:   movie["is_like"].(bool),
			IsUnlike: movie["is_unlike"].(bool),
			VotedAt:  movie["voted_at"].(string),
		})
	}

	// Return the list of voted movies
	return votedMovies, nil
}

// TraceViewership tracks the duration of a movie viewed by a user
func (uc *StatsUseCase) TraceViewership(userID, movieID, duration int) error {
	// Validate if the movie exists
	movieExists, err := uc.StatsRepo.MovieExists(movieID)
	if err != nil {
		return fmt.Errorf("failed to validate movie existence: %w", err)
	}
	if !movieExists {
		return fmt.Errorf("movie not found")
	}

	// Validate if the user has viewed the movie
	movieViewExists, err := uc.StatsRepo.MovieViewExists(userID, movieID)
	if err != nil {
		return fmt.Errorf("failed to validate movie view existence: %w", err)
	}
	if !movieViewExists {
		return fmt.Errorf("user has not viewed this movie")
	}

	// Update the viewing duration
	err = uc.StatsRepo.UpdateViewingDuration(userID, movieID, duration)
	if err != nil {
		return fmt.Errorf("failed to update viewing duration: %w", err)
	}

	return nil
}
