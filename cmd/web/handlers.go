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

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "about.tmpl.html", data)
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
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")

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

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

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

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "Name field cannot be empty")
	form.CheckField(validator.NotBlank(form.Email), "email", "Email field cannot be empty")
	form.CheckField(validator.ValidEmail(form.Email), "email", "Not a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password field cannot be empty")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password field cannot be less than 8 characters")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.ValidEmail(form.Email), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	if app.sessionManager.Exists(r.Context(), "postLoginRedirectURL") {
		url := app.sessionManager.Pop(r.Context(), "postLoginRedirectURL").(string)
		http.Redirect(w, r, url, http.StatusSeeOther)

	} else {
		http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
	}
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// chagne session id again
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	user, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		} else {
			app.serverError(w, err)
		}
	}

	data := app.newTemplateData(r)
	data.User = user
	app.render(w, http.StatusOK, "account.tmpl.html", data)
}

func (app *application) changePassword(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "changePassword.tmpl.html", data)
}

type passwordChangeForm struct {
	CurrentPassword     string `form:"current_pass"`
	NewPassword         string `form:"new_pass"`
	ConfirmNewPassword  string `form:"confirm_new_pass"`
	validator.Validator `form:"-"`
}

func (app *application) changePasswordPost(w http.ResponseWriter, r *http.Request) {
	var form passwordChangeForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPass", "This field cannot be blank")
	// So with the generate hash function, it has a max of 72 bytes for the input, so we need
	// to ensure that it's less than that. We can set a char count on it, but we also would
	// want to make sure that it's alphanumeric (no crazy multi unicode point graphemes)
	// Should unicode points beyond 2 bytes even be considered?
	form.CheckField(
		validator.MinChars(form.NewPassword, 8) && validator.MaxChars(form.NewPassword, 15),
		"newPass",
		"Password must be between 8 and 15 characters long",
	)
	form.CheckField(form.NewPassword == form.ConfirmNewPassword, "confirmNewPass", "Passwords do not match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "changePassword.tmpl.html", data)
		return
	}

	id := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	user, err := app.users.Get(id)
	if err != nil {
		// not sure what the problem would be if the session has an invalid
		// authenticatedUserID since this route got past the Authenticate middleware
		app.serverError(w, err)
		return
	}

	id, err = app.users.Authenticate(user.Email, form.CurrentPassword)

	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("current_pass", "Invalid current password")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "changePassword.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// The user is now authorized to make a password change
	err = app.users.UpdatePassword(id, form.NewPassword)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Password successfully updated")
	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
