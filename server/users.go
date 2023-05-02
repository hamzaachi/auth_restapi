package server

import (
	"net/http"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"hamza.achi/auth/models"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)

	if err := dbInstance.CreateUser(user); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}

	if err := render.Render(w, r, user); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
