package model

// Users represents a user with their email, password, and other details.
type Users struct {
	Id       int    `json:"id"`       // User ID
	Email    string `json:"email"`    // User's Email
	Password string `json:"password"` // User's Password
	Gender   string `json:"gender"`   // User's Gender
	Role     string `json:"role"`     // User's Role
}
