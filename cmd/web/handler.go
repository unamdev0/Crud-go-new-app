package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
	"github.com/unamdev0/go-crud-app/models"
)

func (a *application) homeHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		a.serverError(w, err)
		return
	}

	filter := models.Filter{
		Query:    r.URL.Query().Get("q"),
		Page:     a.readIntDefault(r, "page", 1),
		PageSize: a.readIntDefault(r, "page", 5),
		OrderBy:  r.URL.Query().Get("order_by"),
	}

	posts, meta, err := a.Models.Posts.GetAll(filter)
	if err != nil {
		a.serverError(w, err)
		return
	}

	queryUrl := fmt.Sprintf("page_size=%d&order_by=%s&q=%s", meta.PageSize, filter.OrderBy, filter.Query)
	nextUrl := fmt.Sprintf("%s&page=%d", queryUrl, meta.NextPage)
	prevUrl := fmt.Sprintf("%s&page=%d", queryUrl, meta.PrevPage)

	vars := make(jet.VarMap)
	vars.Set("posts", posts)
	vars.Set("meta", meta)
	vars.Set("nextUrl", nextUrl)
	vars.Set("prevUrl", prevUrl)
	// vars.Set("form", forms.New(r.Form))

	err = a.render(w, r, "index", vars)

	if err != nil {
		log.Fatal(err)
	}

}

func (a *application) commentHandler(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)

	postId, err := strconv.Atoi(chi.URLParam(r, "postID"))

	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	post, err := a.Models.Posts.GetById(postId)
	if err != nil {
		a.serverError(w, err)
		return
	}

	comments, err := a.Models.Comments.GetForPost(postId)
	if err != nil {
		a.serverError(w, err)
		return
	}

	vars.Set("post", post)
	vars.Set("comments", comments)

	err = a.render(w, r, "comments", vars)

	if err != nil {
		log.Fatal(err)
	}

}
