package main

import (
	"net/http"
)

type loginForm struct {
	Username string
	Password string
}

func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = loginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	// TODO Pass and login check here

	token, err := app.generateJWT("akulbek")
	if err != nil {
		app.serverError(w, err)
		return
	}
	cookie := &http.Cookie{Name: "auth_token", Value: token, HttpOnly: true}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusAccepted)
}
