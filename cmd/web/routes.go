package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(a.LoadSession)

	mux.Get("/", a.homeHandler)
	mux.Get("/comment/{postID}", a.commentHandler)

	mux.Get("/login", a.loginHandler)
	mux.Post("/login", a.loginPostHandler)
	mux.Get("/signup", a.signUpHandler)
	mux.Post("/signup", a.signPostUpHandler)
	mux.Get("/logout", a.authRequired(a.logoutHandler))

	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/*", http.StripPrefix("/public", fileServer))

	//Register routes

	return mux
}
