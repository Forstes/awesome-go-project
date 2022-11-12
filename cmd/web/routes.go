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

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/login", app.loginForm)
	router.HandlerFunc(http.MethodPost, "/login", app.loginPost)
	router.HandlerFunc(http.MethodGet, "/signup", app.signupForm)
	router.HandlerFunc(http.MethodPost, "/signup", app.signupPost)
	router.HandlerFunc(http.MethodGet, "/logout", app.logout)
	router.HandlerFunc(http.MethodGet, "/picture/view/:id", app.pictureView)
	router.Handler(http.MethodGet, "/picture/create", app.verifyJWT(http.HandlerFunc(app.pictureUploadForm)))
	router.HandlerFunc(http.MethodPost, "/picture/create", app.pictureUploadPost)
	router.HandlerFunc(http.MethodGet, "/pictures/:path", app.pictureStorage)

	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
