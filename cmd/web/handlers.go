package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"awesome.forstes.go/internal/models"
	"awesome.forstes.go/internal/validator"
	"github.com/gabriel-vasile/mimetype"
	"github.com/julienschmidt/httprouter"
)

func (app *application) pictureStorage(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	path := params.ByName("path")
	http.Redirect(w, r, "http://localhost:9000"+"/pictures/"+path, http.StatusPermanentRedirect)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	pictures, err := app.pictures.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Pictures = pictures

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) pictureView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	picture, err := app.pictures.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	user, err := app.user.GetUserByPicture(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)

		}
		return
	}
	picture.Owner = user.Name
	data := app.newTemplateData(r)
	data.Picture = picture

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) pictureUploadForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

type snippetCreateForm struct {
	Title   string
	Expires int
	validator.Validator
}

func (app *application) pictureUploadPost(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("auth_token")
	claims, err := app.extractClaims(cookie)
	if err != nil {
		app.serverError(w, err)
		return
	}
	userIdStr := claims["user"].(string)
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = r.ParseMultipartForm(1024 * 1024 * 16) // max 16 MB
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	println(userId)
	/*expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	} */

	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Expires: 365,
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	defer file.Close()
	println(userId)
	mtype, err := mimetype.DetectReader(file)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	println(userId)
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")
	form.CheckField(validator.PermittedFileFormat(mtype.String(), "image/png", "image/jpeg"), "file", "File must be a png/jpeg image")
	form.CheckField(validator.MaxMb(handler.Size, 10), "file", "File size should not exceed 10 mb")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	path := app.generatePath(handler.Filename, strconv.Itoa(userId))

	id, err := app.pictures.Insert(userId, form.Title, path, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = app.objStorage.UploadObject("pictures", path, file, handler.Size, mtype.String())
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/picture/view/%d", id), http.StatusSeeOther)
}
