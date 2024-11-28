package model

// Movies represents a movie record.
type Movies struct {
	ID          int    `json:"id"`                 // Movie ID
	Title       string `json:"title"`              // Movie title
	Description string `json:"description"`        // Movie description
	Duration    string `json:"duration"`           // Movie duration
	Artist      string `json:"artist"`             // Main artist (actor/director)
	GenreID     *int   `json:"genre_id,omitempty"` // Genre ID (nullable)
	WatchURL    string `json:"watch_url"`          // URL to watch the movie
}
