package model

// MovieStatsView represents statistics for a movie.
type MovieStatsView struct {
	ID    int    `json:"id"`    // Movie ID
	Title string `json:"title"` // Movie Title
	Views int    `json:"views"` // Number of Views
}

// MovieVotesView represents vote statistics for a movie.
type MovieStatsVote struct {
	ID        int    `json:"id"`         // Movie ID
	Title     string `json:"title"`      // Movie Title
	VoteCount int    `json:"vote_count"` // Number of Vote
}

// GenreStats represents statistics for a genre.
type GenreStats struct {
	ID    int    `json:"id"`    // Genre ID
	Name  string `json:"name"`  // Genre Name
	Views int    `json:"views"` // Number of Views
}

// StatsModel summarizes movie and genre statistics.
type StatsModelView struct {
	MostViewedMovie []MovieStatsView `json:"most_viewed_movie"` // Most Viewed Movie
	MostViewedGenre []GenreStats     `json:"most_viewed_genre"` // Most Viewed Genre
}

type StatsModelVote struct {
	MostVotedMovie  []MovieStatsVote `json:"most_voted_movie"`  // Most Vote Movie
	MostViewedGenre []GenreStats     `json:"most_viewed_genre"` // Most Viewed Genre
}

// MovieView represents a record of a movie viewed by a user.
type MovieView struct {
	MovieID  int    `json:"movie_id"`  // Movie ID
	UserID   int    `json:"user_id"`   // User ID
	ViewedAt string `json:"viewed_at"` // Timestamp of the view
}

// UserVotedMovie represents a user's vote on a movie.
type UserVotedMovie struct {
	MovieID  int    `json:"movie_id"`  // Movie ID
	Title    string `json:"title"`     // Movie Title
	IsLike   bool   `json:"is_like"`   // Whether the user likes the movie
	IsUnlike bool   `json:"is_unlike"` // Whether the user dislikes the movie
	VotedAt  string `json:"voted_at"`  // Timestamp of the vote
}
