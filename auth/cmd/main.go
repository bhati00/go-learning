package main

import (
	"log"
	"net/http"

	"github.com/bhati00/go-learning/auth/internal/handler"
	"github.com/bhati00/go-learning/auth/internal/middleware"
	"github.com/bhati00/go-learning/auth/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func main() {
	tokenAuth = jwtauth.New("HS256", []byte("secret-key"), nil)

	db, err := repository.InitDB("internal/data/auth.db")
	if err != nil {
		log.Fatalf("%v", err)
	}
	handler.Init(db)
	r := chi.NewRouter()
	r.Post("/api/register", handler.RegisterHandler)
	r.Post("/api/login", handler.LoginHandler)

	r.Group(func(protected chi.Router) {
		protected.Use(jwtauth.Verifier(tokenAuth))
		protected.Use(jwtauth.Authenticator(tokenAuth))
		protected.Use(middleware.RequiredAuth)
		protected.Get("/api/profile", handler.ProfileHandler)
	})

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
