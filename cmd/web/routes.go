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
	//router.HandlerFunc(http.MethodGet, "/login", app.loginForm)
	//router.HandlerFunc(http.MethodPost, "/login", app.loginPost)
	//router.HandlerFunc(http.MethodGet, "/picture/view/:id", app.pictureView)
	router.Handler(http.MethodGet, "/picture/create", app.verifyJWT(http.HandlerFunc(app.pictureUploadForm)))
	//router.HandlerFunc(http.MethodPost, "/picture/create", app.pictureUploadPost)
	router.HandlerFunc(http.MethodGet, "/pictures/:path", app.pictureStorage)

	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
