package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.cozycole.net/internal/models"
)

// Define a templateData type to act as the holding
// struct for any dynamic data we want to pass to the
// HTML templates. Since the ExecuteTemplate only accepts one
// struct for inserting data and data can come from many sources,
// you need to combine it all into one
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Init global variable which maps string func names to
// functions to be used within templates (since you can call
// functions from template). NOTE: The tempalte functions should only
// return a single value
var functions = template.FuncMap{
	"humanDate": humanDate,
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

		// Register the funcMap before parsing the files
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")

		if err != nil {
			return nil, err
		}

		// Add all the template partials that could be inserted into
		// the given page template
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
