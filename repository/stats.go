package repository

import (
	"database/sql"
	"fmt"
	"movies/model"
)

type StatsRepo struct {
	DB *sql.DB
}

func NewStatsRepo(DB *sql.DB) *StatsRepo {
	return &StatsRepo{DB: DB}
}

// GetMostViewedMovies retrieves all movies with the highest number of views from the database.
func (r *StatsRepo) GetMostViewedMovies() ([]model.MovieStatsView, error) {
	query := `
		SELECT m.id, m.title, COUNT(mv.movie_id) AS views
		FROM movies m
		LEFT JOIN movie_views mv ON m.id = mv.movie_id
		GROUP BY m.id, m.title
		ORDER BY views DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []model.MovieStatsView
	var maxViews int

	for rows.Next() {
		var movie model.MovieStatsView

		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Views); err != nil {
			return nil, err
		}

		// Stop adding movies if their views are less than the maxViews
		if len(movies) > 0 && movie.Views < maxViews {
			break
		}

		// Set maxViews for the first iteration
		if len(movies) == 0 {
			maxViews = movie.Views
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

// GetMostViewedGenres retrieves all genres with the highest number of views from the database.
func (r *StatsRepo) GetMostViewedGenres() ([]model.GenreStats, error) {
	query := `
		SELECT g.id, g.name, COUNT(mv.movie_id) AS views
		FROM genres g
		LEFT JOIN movies m ON g.id = m.genre_id
		LEFT JOIN movie_views mv ON m.id = mv.movie_id
		GROUP BY g.id, g.name
		ORDER BY views DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []model.GenreStats
	var maxViews int

	for rows.Next() {
		var genre model.GenreStats

		if err := rows.Scan(&genre.ID, &genre.Name, &genre.Views); err != nil {
			return nil, err
		}

		// Stop adding genres if their views are less than the maxViews
		if len(genres) > 0 && genre.Views < maxViews {
			break
		}

		// Set maxViews for the first iteration
		if len(genres) == 0 {
			maxViews = genre.Views
		}

		genres = append(genres, genre)
	}

	return genres, nil
}

// GetMostVotedMovies retrieves all movies with the most positive votes (is_like = 1).
func (repo *StatsRepo) GetMostVotedMovies() ([]model.MovieStatsVote, error) {
	query := `
		SELECT m.id, m.title, COUNT(uv.id) AS vote_count
		FROM movies m
		LEFT JOIN user_votes uv ON m.id = uv.movie_id AND uv.is_like = 1
		GROUP BY m.id, m.title
		ORDER BY vote_count DESC
	`

	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []model.MovieStatsVote
	var maxVotes int

	for rows.Next() {
		var movie model.MovieStatsVote
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.VoteCount); err != nil {
			return nil, err
		}

		// Stop adding movies if their vote count is less than the maxVotes
		if len(movies) > 0 && movie.VoteCount < maxVotes {
			break
		}

		// Set maxVotes for the first iteration
		if len(movies) == 0 {
			maxVotes = movie.VoteCount
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

// CheckVote checks if the user has liked a particular movie.
func (r *StatsRepo) CheckVote(userID, movieID int) (bool, error) {
	// SQL query to check if the user has liked the movie
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM user_votes 
			WHERE user_id = ? AND movie_id = ? AND is_like = ?
		)
	`
	var exists bool
	// Execute the query and check if the user has liked the movie
	err := r.DB.QueryRow(query, userID, movieID, 1).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CheckUnVote checks if the user has disliked a particular movie.
func (r *StatsRepo) CheckUnVote(userID, movieID int) (bool, error) {
	// SQL query to check if the user has disliked the movie
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM user_votes 
			WHERE user_id = ? AND movie_id = ? AND is_unlike = ?
		)
	`
	var exists bool
	// Execute the query and check if the user has disliked the movie
	err := r.DB.QueryRow(query, userID, movieID, 1).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// AddVote adds a "like" vote for a movie by the user.
func (r *StatsRepo) AddVote(userID, movieID int) error {
	// SQL query to insert or update the "like" vote for a movie
	query := `
		INSERT INTO user_votes (user_id, movie_id, created_at, is_like, is_unlike)
		VALUES (?, ?, NOW(), 1, 0)
		ON DUPLICATE KEY UPDATE is_like = 1, is_unlike = 0
	`
	// Execute the query to add a "like" vote
	_, err := r.DB.Exec(query, userID, movieID)
	return err
}

// AddUnVote adds a "dislike" vote for a movie by the user.
func (r *StatsRepo) AddUnVote(userID, movieID int) error {
	// SQL query to insert or update the "dislike" vote for a movie
	query := `
		INSERT INTO user_votes (user_id, movie_id, created_at, is_like, is_unlike)
		VALUES (?, ?, NOW(), 0, 1)
		ON DUPLICATE KEY UPDATE is_like = 0, is_unlike = 1
	`
	// Execute the query to add a "dislike" vote
	_, err := r.DB.Exec(query, userID, movieID)
	return err
}

// GetUserVotedMovies retrieves all movies that the user has voted on.
func (r *StatsRepo) GetUserVotedMovies(userID int) ([]map[string]interface{}, error) {
	// SQL query to get all movies voted by the user
	query := `
		SELECT 
			m.id AS movie_id, 
			m.title, 
			uv.is_like, 
			uv.is_unlike, 
			uv.created_at AS voted_at
		FROM movies m
		INNER JOIN user_votes uv ON m.id = uv.movie_id
		WHERE uv.user_id = ?
		ORDER BY uv.created_at DESC
	`

	// Execute the query to get all the user's voted movies
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Store the list of voted movies
	var votedMovies []map[string]interface{}
	for rows.Next() {
		var movieID int
		var title string
		var isLike, isUnlike bool
		var votedAt string

		// Scan the row into variables
		err = rows.Scan(&movieID, &title, &isLike, &isUnlike, &votedAt)
		if err != nil {
			return nil, err
		}

		// Append the movie details to the list
		votedMovies = append(votedMovies, map[string]interface{}{
			"movie_id":  movieID,
			"title":     title,
			"is_like":   isLike,
			"is_unlike": isUnlike,
			"voted_at":  votedAt,
		})
	}

	// Check for errors during row iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return the list of voted movies
	return votedMovies, nil
}

// MovieExists checks if a movie with the given ID exists.
func (r *StatsRepo) MovieExists(movieID int) (bool, error) {
	// SQL query to check if the movie exists in the database
	query := `
		SELECT COUNT(1)
		FROM movies
		WHERE id = ?
	`

	var count int
	// Execute the query to check for movie existence
	err := r.DB.QueryRow(query, movieID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check movie existence: %w", err)
	}

	// Return true if the movie exists
	return count > 0, nil
}

// MovieViewExists checks if a movie has views recorded in the movie_views table.
func (r *StatsRepo) MovieViewExists(userID int, movieID int) (bool, error) {
	// SQL query to check if the movie has views in the database
	query := `
		SELECT COUNT(1)
		FROM movie_views
		WHERE user_id = ? AND movie_id = ?
	`
	var count int
	// Execute the query to check if movie views exist
	err := r.DB.QueryRow(query, userID, movieID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if movie_id exists in movie_views: %w", err)
	}

	// Return true if there are views recorded for the movie
	return count > 0, nil
}
