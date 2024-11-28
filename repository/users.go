package repository

import (
	"database/sql"
	"log"
	"movies/model"
)

type UsersRepo struct {
	DB *sql.DB
}

func NewUsersRepo(DB *sql.DB) *UsersRepo {
	return &UsersRepo{DB: DB}
}

// Create inserts a new user into the users table.
func (ur *UsersRepo) Create(user *model.Users) (*model.Users, error) {
	// SQL query to insert a new user record into the database.
	sql_insert := "INSERT INTO users (email, password, gender, role) VALUES (?,?,?,?)"

	// Execute the insert query with the provided user details.
	_, err := ur.DB.Exec(sql_insert, user.Email, user.Password, user.Gender, user.Role)

	// If there was an error executing the query, log it and return the error.
	if err != nil {
		log.Println("ERR : ", err)
		return nil, err
	}

	// Return the user object if insertion was successful.
	return user, nil
}

// IsEmailExists checks if a user with the provided email already exists in the database.
func (ur *UsersRepo) IsEmailExists(email string) (bool, error) {
	// SQL query to check if the email exists in the users table.
	sql_query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"

	// Variable to store the result of the query.
	var exists bool
	// Execute the query and scan the result into the exists variable.
	err := ur.DB.QueryRow(sql_query, email).Scan(&exists)

	// If there's an error checking the email existence, log it and return the error.
	if err != nil {
		log.Println("Error checking email existence: ", err)
		return false, err
	}

	// Return the result indicating whether the email exists or not.
	return exists, nil
}

// GetUser retrieves a user record based on the provided email.
func (ur *UsersRepo) GetUser(user *model.Users) (*model.Users, error) {
	// SQL query to retrieve a user based on the email.
	sql_query := "SELECT id, email, password, gender, role FROM users WHERE email = ?"
	var result model.Users

	// Execute the query to get the user data.
	row := ur.DB.QueryRow(sql_query, user.Email)

	// Scan the result into the result model.
	if err := row.Scan(&result.Id, &result.Email, &result.Password, &result.Gender, &result.Role); err != nil {
		// If no rows were found (email does not exist), return an ErrNoRows error.
		if err == sql.ErrNoRows {
			log.Println("Email not found: ", user.Email)
			return nil, sql.ErrNoRows
		}

		// Handle any other errors that might occur.
		log.Println("ERR: ", err)
		return nil, err
	}

	// Return the user if found.
	return &result, nil
}
