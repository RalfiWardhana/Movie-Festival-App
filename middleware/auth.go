package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the structure of the JWT claims
type Claims struct {
	UserID int    `json:"id"`   // User ID in the JWT claim
	Role   string `json:"role"` // User role in the JWT claim
	jwt.RegisteredClaims
}

// AuthMiddleware authenticates requests using JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	var jwtSecret = []byte("your_secret_key") // Secret key for JWT
	return func(c *gin.Context) {
		// Retrieve the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// Parse the token format (e.g., "Bearer <token>")
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString = parts[1]

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store claims in the context for further use
		c.Set("claims", claims)
		c.Next()
	}
}

// RoleMiddleware restricts access to users with specific roles
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve claims from the context
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Missing claims"})
			c.Abort()
			return
		}

		// Assert claims to the expected type
		userClaims, ok := claims.(*Claims)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid claims"})
			c.Abort()
			return
		}

		// Check if the user's role is in the allowed roles
		for _, role := range allowedRoles {
			if userClaims.Role == role {
				c.Next()
				return
			}
		}

		// If no match, respond with forbidden
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		c.Abort()
	}
}
