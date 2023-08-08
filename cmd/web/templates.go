package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.cozycole.net/internal/models"
)

// Define a templateData type to act as the holding
// struct for any dynamic data we want to pass to the
// HTML templates. Since the ExecuteTemplate only accepts one
// struct for inserting data and data can come from many sources,
// you need to combine it all into one
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

// Getting mapping of html page filename to template set for the page
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
