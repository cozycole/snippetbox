package main

import "snippetbox.cozycole.net/internal/models"

// Define a templateData type to act as the holding
// struct for any dynamic data we want to pass to the
// HTML templates. Since the ExecuteTemplate only accepts one
// struct for inserting data and data can come from many sources,
// you need to combine it all into one
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
