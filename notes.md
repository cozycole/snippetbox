# Creating a module
Module path: a standard name or identifier for the project. Almost any string can be used but it needs to be unique to avoid import conflicts with other projects. Common convention is to base a module path on a URL you own. Example: snippetbox.cozycole.net

You then initialize the go module with that name 

```bash
go mod init snippetbox.cozycole.net
```

This generates a go.mod file which signals the directory it is stored in is a go module.

# Go's servemux (part of the http package)

servemux is Go's url router mapping paths to handlers.
servemux provides two types of URL patterns: fixed and subtree paths
subtree paths end with a slash (ex: / or /static/) and matches all urls with that root
fixed ones (ex: /snippet/view) only match the specific url

Longer paths take precedence!
Requests to paths missing trailing slash get automatically redirected ( so if /foo/ is registered, /foo -> /foo/)

Go's servemux lacks the ability to route based on method or regexp-based patterns. You'll need a third party option for that

# Customizing HTTP Headers

We normally pass the http.ResponseWriter to another function which sends a response for us
instead of writing w.Write() and w.WriteHeader()

It's best practice to use the net/http constants for methods and response codes (ex: http.MethodPost and http.StatusMethodNotAllowed)

Go automatically sets three headers for you: Date, Content-Length and Content-Type

The Content-Type is determined automatically (with http.DetectContentType()) to guess the content type.
NOTE: it detects JSON as plain text so you need to set JSON manually.

Headers can be manipulated with w.Header().Set,.Add,.Del,.Get or by editing the header directly with w.Header()["Key"] = []string{"val"}

# Project Structure

Refer to for best practices: https://peter.bourgon.org/go-best-practices-2016/#repository-structure

## This project's structure

- cmd: Contains application specific code for executable applications in the project.
- internal: contains non-application-specific code used in the project (ex: validation helpers and SQL db models)
- ui: user-interface assets (html, CSS, images, javascript)

# Templates

Go provides templates which are used to break down html documents into dynamic componenets. You define a template like:

```html
{{define "template-name"}} <p> This is my element </p> {{end}}
```

and you reference it in code with:
```
{{template-name}}
```

# What is a handler?

A handler is object that implements the following type interface:

```go
type Hanlder interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

So in its simplest form, we he could create the following handler:

```go
type home struct {}

func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is my home page"))
}

// Then register it with servemux
mux := http.NewServerMux()
mux.Handle("/", &home{})
```

This is kinda clunky tho since we don't need to make an object just for that. We instead create a function

```go
func home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Home page"))
}

mux := http.NewServerMux()
mux.HandleFunc("/", home)
```

The method HandleFunc adds a ServeHTTP method to the function object then registers the handler.

# So what happens when a request is made to the server?

Since the servemux object itself satsifies the handler interface (has a ServeHTTP method), when the server
receieves a request, it calls the sermux's ServeHTTP() method. Servemux's ServeHTTP method looks up the relevant 
handler object based on the request URL and then calls that handler's ServeHTTP() method.

**NOTE**
Requests are handled in parallel so you need to account for race conditions when accessing shared resources.

