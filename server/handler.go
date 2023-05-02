package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"hamza.achi/auth/db"
)

var dbInstance db.Database

func NewRouter(db db.Database) chi.Router {
	r := chi.NewRouter()
	dbInstance = db

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcom Home!"))
	})
	r.Post("/", CreateUser)
	r.Post("/login", Login)
	return r
}
