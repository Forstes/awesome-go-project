package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"awesome.forstes.go/internal/models"
	"awesome.forstes.go/internal/validator"
)

type loginForm struct {
	Name                string
	Password            string
	validator.Validator `form:"-"`
}

type signupForm struct {
	Name                string `form:"name"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = loginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	// TODO Pass and login check here
	form := loginForm{
		Name:     r.PostForm.Get("name"),
		Password: r.PostForm.Get("password"),
	}
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	id, err := app.user.Authenticate(form.Name, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Name or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	token, err := app.generateJWT(fmt.Sprint(id))
	if err != nil {
		app.serverError(w, err)
		return
	}
	cookie := &http.Cookie{Name: "auth_token", Value: token, Expires: time.Now().Add(8 * time.Hour), HttpOnly: true}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = signupForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	form := signupForm{
		Name:     r.PostForm.Get("name"),
		Password: r.PostForm.Get("password"),
	}
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	id, err := app.user.Insert(form.Name, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateName) {
			form.AddFieldError("name", "this name is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}

	}

	if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Fprintf(w, "Create a new user with ID: %v", id)

}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {

}
