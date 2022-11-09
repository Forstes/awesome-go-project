package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	router.HandlerFunc(http.MethodGet, "/", app.home)
	//router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.pictureView)
	//router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	//router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)
	router.HandlerFunc(http.MethodGet, "/pictures/:path", app.pictureStorage)

	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
