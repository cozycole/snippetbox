package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.cozycole.net/internal/models"
	"snippetbox.cozycole.net/internal/validator"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

// Include struct tags which tell the decoder how to map HTML form values
// into the different struct field. For example, here we're telling the decoder
// to store the value from the HTML form input with the name "title" in the Title field. The struct tag `form:"-`
// tells the decoder to completely ignore a field during decoding.
type snippetCreateForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	// adds the validator package as an attribute
	// meaning public functions of validator.Validator
	// act as methods
	validator.Validator `form:"-"`
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	// Loads the values from the sent Form into the snippetCreateForm based
	// on matching `struct-tags`
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Since the Validator type is embedded in the snippetCreateForm, we can
	// call CheckField directly on the object.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")

	if !form.Valid() {
		// sending a new html form with errors if it's not valid
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// We need to initialize the form data since the template needs it to
	// render. It's a good place to put default values for the fields too (e.g. Expires = 365 will default that option in the template)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl.html", data)
}
