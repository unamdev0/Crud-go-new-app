package main

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
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

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {

		a.session.Put(r.Context(), "Key", "Value")
		err := a.render(w, r, "index", nil)

		if err != nil {
			log.Println("dsfafdsa")
			log.Fatal(err)
		}

	})

	mux.Get("/comments", func(w http.ResponseWriter, r *http.Request) {

		vars := make(jet.VarMap)
		vars.Set("Key", a.session.GetString(r.Context(), "Key"))
		err := a.render(w, r, "index", vars)

		if err != nil {
			log.Fatal(err)
		}

	})
	return mux
}
