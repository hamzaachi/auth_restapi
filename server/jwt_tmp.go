package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type loginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

var (
	secretKey       = []byte("your-secret-key") // Replace with your secret key
	refreshTokenMap = make(map[string]uuid.UUID)
)

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/refresh", refreshHandler)
	http.HandleFunc("/protected", protectedHandler)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Simulating successful authentication
	userID := 123
	accessToken, err := generateAccessToken(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken, err := generateRefreshToken(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expiresIn := 3600 // Token expiration time in seconds

	response := loginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify the refresh token
	refreshToken := req.RefreshToken
	if !isValidRefreshToken(refreshToken) {
		http.Error(w, "Invalid refresh token", http.StatusBadRequest)
		return
	}

	// Extract the user ID from the refresh token (You may store additional user data in the token itself)
	userID, err := getUserIDFromRefreshToken(refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate a new access token
	accessToken, err := generateAccessToken(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the new access token in the response
	response := loginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	// Verify the access token
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
}
// Access token is valid, serve the protected content
fmt.Fprintln(w, "Protected content!")
}

func generateAccessToken(userID int) (string, error) {
token := jwt.New(jwt.SigningMethodHS256)

claims := token.Claims.(jwt.MapClaims)
claims["userID"] = userID
claims["exp"] = time.Now().Add(time.Hour).Unix()

tokenString, err := token.SignedString(secretKey)
if err != nil {
	return "", err
}

return tokenString, nil

}

func generateRefreshToken(userID int) (string, error) {
refreshToken := uuid.New().String()
// Store the refresh token in memory or persistent storage
refreshTokenMap[refreshToken] = uuid.FromInt(int64(userID))

return refreshToken, nil
}

func isValidRefreshToken(refreshToken string) bool {
_, exists := refreshTokenMap[refreshToken]
return exists
}

func getUserIDFromRefreshToken(refreshToken string) (int, error) {
userUUID, exists := refreshTokenMap[refreshToken]
if !exists {
return 0, fmt.Errorf("invalid refresh token")
}
return int(userUUID.Int64()), nil
}
