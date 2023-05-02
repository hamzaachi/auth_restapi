package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"hamza.achi/auth/models"
)

var mySigningKey = []byte("NHazJ-mcQ9RudYmJUO922jqFzBYOzlworO6PvJddPaI") // Replace with your own secret key

func authenticate(w http.ResponseWriter, r *http.Request) {
	// Get the token from the Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		response := models.Response{
			Status:  http.StatusUnauthorized,
			Message: "No token provided",
		}
		sendResponse(w, http.StatusUnauthorized, response)
		return
	}

	// Verify the token
	valid, username := verifyToken(tokenString)
	if valid {
		response := models.Response{
			Status:  http.StatusOK,
			Message: fmt.Sprintf("Authentication successful for user: %s", username),
		}
		sendResponse(w, http.StatusOK, response)
	} else {
		response := models.Response{
			Status:  http.StatusUnauthorized,
			Message: "Invalid token",
		}
		sendResponse(w, http.StatusUnauthorized, response)
	}
}

func verifyToken(tokenString string) (bool, string) {
	// Remove the "Bearer " prefix from the token string
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method used
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return mySigningKey, nil
	})

	if err != nil {
		return false, ""
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Token is valid
		username := claims["username"].(string)
		return true, username
	}

	return false, ""
}

// Rest of the code...

func generateToken(username string) (string, error) {
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Set the token expiration time (1 day)
	})

	// Generate the token string
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", fmt.Errorf("Failed to generate token: %v", err)
	}

	return tokenString, nil
}

func comparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func sendResponse(w http.ResponseWriter, statusCode int, response models.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
