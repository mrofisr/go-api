package middleware

import (
	"net/http"
	jwt "github.com/golang-jwt/jwt/v5"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract JWT token from the request
		tokenString := extractToken(r)

		// Validate JWT token
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Implement your secret key here

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func extractToken(r *http.Request) string {
	// Extract token from Authorization Bearer header or from query parameter or from cookie
}
