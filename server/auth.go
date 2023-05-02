package server

import (
	"net/http"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"hamza.achi/auth/models"
)

func Login(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}

	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	u, err := dbInstance.GetUserById(user.Username)
	if err != nil {
		render.Render(w, r, ErrUnauthorized)
		return
	}

	hashedPassword := u.Password

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		render.Render(w, r, ErrUnauthorized)
		return
	}

	tokenString, err := generateToken(user.Username)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
	}

	response := models.Response{
		Status:  http.StatusOK,
		Message: "Authentication successful",
		Token:   tokenString,
	}
	render.Render(w, r, &response)
	return
	//w.Write([]byte("Success"))
}
