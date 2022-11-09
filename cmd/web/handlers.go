package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"awesome.forstes.go/internal/models"
	"github.com/gabriel-vasile/mimetype"
	"github.com/julienschmidt/httprouter"
)

// BAD
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
	/*
		// TODO someone clean this code by moving to helpers, please
		var pics []*dto.PictureResponse
		for _, p := range pictures {
			pics = append(pics, &dto.PictureResponse{PostedBy: p.Owner, Title: p.Title, Created: p.Created, Expires: p.Expires})
			app.infoLog.Println(p.Path)
			url, err := app.objStorage.GetObjectPresigned("pictures", "ranalda.png")
			if err != nil {
				app.serverError(w, err)
				return
			}
			pics[len(pics)-1].Url = url
		}
	*/
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
	/*
		bytes, err := app.objStorage.GetObject("pictures", picture.Path)
		if err != nil {
			app.serverError(w, err)
			return
		}
	*/
	data := app.newTemplateData(r)
	data.Picture = picture
	//&dto.PictureResponse{PostedBy: picture.Owner, Title: picture.Title, Created: picture.Created, Expires: picture.Expires}

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	//validator.Validator
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	// 16 mb
	r.ParseMultipartForm(5 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	defer file.Close()

	mtype, err := mimetype.DetectReader(file)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// TODO Do some validation shit
	/* 	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	   	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	   	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	   	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	   	if !form.Valid() {
	   		data := app.newTemplateData(r)
	   		data.Form = form
	   		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
	   		return
	   	}*/

	id, err := app.pictures.Insert(form.Title, form.Content, form.Expires)
	app.objStorage.UploadObject("pictures", handler.Filename, file, handler.Size, mtype.String())

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
