package repository

import (
	"database/sql"
	"fmt"
	"log"
	"movies/model"
)

type MoviesRepo struct {
	DB *sql.DB
}

func NewMoviesRepo(DB *sql.DB) *MoviesRepo {
	return &MoviesRepo{DB: DB}
}

// Create inserts a new movie into the database.
func (pr *MoviesRepo) Create(Movie *model.Movies) (*model.Movies, error) {
	sql_insert := "INSERT INTO movies (title, description, duration, artist, genre_id, watch_url) VALUES (?,?,?,?,?,?)"
	_, err := pr.DB.Exec(sql_insert, Movie.Title, Movie.Description, Movie.Duration, Movie.Artist, Movie.GenreID, Movie.WatchURL)

	if err != nil {
		log.Println("ERR insert Movies : ", err)
		return nil, err
	}

	return Movie, nil
}

// Update updates an existing movie record in the database.
func (repo *MoviesRepo) Update(movieID int, updates map[string]interface{}) error {
	// Start building the UPDATE query string.
	query := "UPDATE movies SET "
	params := []interface{}{}

	// Add each column to the update query dynamically.
	for column, value := range updates {
		query += fmt.Sprintf("%s = ?, ", column)
		params = append(params, value)
	}

	// Remove the trailing comma and finalize the query.
	query = query[:len(query)-2]
	query += " WHERE id = ?"
	params = append(params, movieID)

	// Execute the query with the provided parameters.
	_, err := repo.DB.Exec(query, params...)
	return err
}

// GetAllMoviesWithPagination retrieves a list of movies from the database with pagination.
func (r *MoviesRepo) GetAllMoviesWithPagination(page, limit int) ([]model.Movies, error) {
	// Calculate the offset based on the page and limit.
	offset := (page - 1) * limit
	query := fmt.Sprintf(`
		SELECT id, title, description, duration, artist, genre_id, watch_url
		FROM movies
		ORDER BY id ASC
		LIMIT %d OFFSET %d
	`, limit, offset)

	// Execute the query and get the result rows.
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the rows into a slice of movie models.
	var movies []model.Movies
	for rows.Next() {
		var movie model.Movies
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.Artist, &movie.GenreID, &movie.WatchURL); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	// Return the list of movies.
	return movies, nil
}

// SearchMovies searches for movies by title, description, artist, or genre ID.
func (r *MoviesRepo) SearchMovies(title string, description string, artist string, genreID int) ([]model.Movies, error) {
	// Start building the search query.
	sqlQuery := "SELECT id, title, description, duration, artist, genre_id, watch_url FROM movies WHERE 1=1"
	var queryParams []interface{}

	// Dynamically add search conditions for each parameter.
	if title != "" {
		sqlQuery += " AND title LIKE ?"
		queryParams = append(queryParams, "%"+title+"%")
	}
	if description != "" {
		sqlQuery += " AND description LIKE ?"
		queryParams = append(queryParams, "%"+description+"%")
	}
	if artist != "" {
		sqlQuery += " AND artist LIKE ?"
		queryParams = append(queryParams, "%"+artist+"%")
	}
	if genreID != 0 {
		sqlQuery += " AND genre_id = ?"
		queryParams = append(queryParams, genreID)
	}

	// Execute the query and get the result rows.
	rows, err := r.DB.Query(sqlQuery, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the rows into a slice of movie models.
	var moviesList []model.Movies
	for rows.Next() {
		var movie model.Movies
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.Artist, &movie.GenreID, &movie.WatchURL); err != nil {
			return nil, err
		}
		moviesList = append(moviesList, movie)
	}

	// Return the list of movies matching the search.
	return moviesList, nil
}

// SaveMovieView saves a record of a user viewing a movie in the movie_views table.
func (r *MoviesRepo) SaveMovieView(view *model.MovieView) error {
	query := `
		INSERT INTO movie_views (movie_id, user_id, viewed_at)
		VALUES (?, ?, ?)
	`

	// Execute the insert query with the view data.
	_, err := r.DB.Exec(query, view.MovieID, view.UserID, view.ViewedAt)
	if err != nil {
		return fmt.Errorf("failed to save movie view: %w", err)
	}

	// Return nil if successful.
	return nil
}

// HasViewed checks if a user has already viewed a movie.
func (r *MoviesRepo) HasViewed(movieID int, userID int) (bool, error) {
	query := `
		SELECT COUNT(1) 
		FROM movie_views 
		WHERE movie_id = ? AND user_id = ?
	`

	// Execute the query and check the count.
	var count int
	err := r.DB.QueryRow(query, movieID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check view status: %w", err)
	}

	// Return true if the user has viewed the movie, else false.
	return count > 0, nil
}

// MovieExists checks if a movie with the given ID exists in the database.
func (r *MoviesRepo) MovieExists(movieID int) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM movies
		WHERE id = ?
	`

	// Execute the query and check if the movie exists.
	var count int
	err := r.DB.QueryRow(query, movieID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check movie existence: %w", err)
	}

	// Return true if the movie exists, else false.
	return count > 0, nil
}
