package main

import (
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"awesome.forstes.go/internal/models"
)

type templateData struct {
	CurrentYear int
	Picture     *models.Picture
	Pictures    []*models.Picture
	Form        any
	Flash       *http.Cookie
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	// Get slice of all the filepaths match the pattern
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the file name (like 'home.tmpl.html') from the full filepath
		name := filepath.Base(page)

		// Create template set from base template
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Add partials to template set
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Add page itself to template set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
