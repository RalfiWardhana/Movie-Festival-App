package handler

import (
	"database/sql"
	"log"
	"movies/model"
	"movies/usecase"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	UsersUsecase *usecase.UsersUseCase
}

func NewUsersHandler(UsersUseCase *usecase.UsersUseCase) *UsersHandler {
	return &UsersHandler{UsersUsecase: UsersUseCase}
}

func (ph *UsersHandler) Register(c *gin.Context) {
	var users model.Users

	// Bind request body to users struct
	err := c.ShouldBind(&users)
	log.Println("Users : ", users)

	if err != nil {
		log.Println("Error binding: ", err)
		c.JSON(400, gin.H{"error": "Bad request", "message": err.Error()})
		return
	}

	// Validate that all required fields are provided
	if users.Email == "" || users.Password == "" || users.Gender == "" || users.Role == "" {
		c.JSON(400, gin.H{"error": "Bad request", "message": "All fields (email, password, gender, role) are required"})
		return
	}

	// Validate the Role field
	if users.Role != "admin" && users.Role != "user" {
		c.JSON(400, gin.H{"error": "Bad request", "message": "Role must be 'admin' or 'user'"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password: ", err)
		c.JSON(500, gin.H{"error": "Internal server error", "message": "Failed to hash password"})
		return
	}
	users.Password = string(hashedPassword)

	// Call the usecase to create the user
	UsersResult, err := ph.UsersUsecase.Create(&users)
	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			c.JSON(400, gin.H{"error": "Bad request", "message": "Email already exists"})
		} else {
			log.Println("Error creating user: ", err)
			c.JSON(500, gin.H{"error": "Internal server error", "message": err.Error()})
		}
		return
	}

	// Respond with success
	c.JSON(201, gin.H{"message": "User registered successfully", "user": UsersResult})
}

func (ph *UsersHandler) Login(c *gin.Context) {
	var user model.Users
	var jwtSecret = []byte("your_secret_key") // Securely store this key

	// JWT claims structure
	type Claims struct {
		UserID int    `json:"id"`   // User ID in the JWT claim
		Role   string `json:"role"` // User role in the JWT claim
		jwt.RegisteredClaims
	}

	// Bind JSON request to user struct
	if err := c.ShouldBind(&user); err != nil {
		log.Println("ERR: ", err)
		c.JSON(400, gin.H{"error": "Bad request", "message": err.Error()})
		return
	}

	log.Println("User from handler: ", user)

	// Get user data from the database
	result, err := ph.UsersUsecase.GetUser(&user)
	if err != nil {
		// Check if the error is due to no rows found
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "Not Found", "message": "Email not registered. Please register first."})
			return
		}

		// Handle other errors
		log.Println("ERR: ", err)
		c.JSON(500, gin.H{"error": "Internal server error", "message": err.Error()})
		return
	}

	// Validate email and password
	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized", "message": "Email or password is invalid"})
		return
	}

	// Create JWT token for successful login
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	claims := &Claims{
		UserID: result.Id,
		Role:   result.Role, // Include role in the token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Sign the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Error signing token: ", err)
		c.JSON(500, gin.H{"error": "Internal server error", "message": err.Error()})
		return
	}

	// Return the JWT token to the client
	c.JSON(200, gin.H{"message": "Login successful", "token": tokenString})
}
