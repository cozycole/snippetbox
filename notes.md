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

# How to configure settings at runtime

For example, if we want to change the port we are running the application on we need to change the string in main.go.

The best way to do this is use the flag package and parse flag variables when executing the application. The best is to pair that
with env variables. For example:

```bash
    export SNIPPETBOX_ADDR=":9999"
    go run ./cmd/web -addr=$SNIPPETBOX_ADDR
```

# Logging

You can create different log objects that with configurable prefixes. For example:

```go
infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
```

The second arg is the prefix and the next is a list of flags for certain values to add (date, time and file name).

Now logging is decoupled, making it simple to manage the logs differently depending on the env. we can log to different
on disk files like so:

```bash
go run ./cmd/web >> /tmp/info.log 2>>/tmp/error.log
```

SIDENOTE: avoid using Panic() and Fatal() outside of your main function.


# Using MySQL

## Create a new user

We don't want to connect as root, but instead as a db user w/ restricted permissions.

```sql
CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
-- Important: Make sure to swap 'pass' with a password of your own choosing.
ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';
```

now the database can't be accidentally or malisciously deleted.

## Install database driver

We need a driver to act as a middle man between the Go applications and MySQL. We will use go-sql-driver/mysql

```bash
go get github.com/go-sql-driver/mysql@v1
```
@v1 indicates we want the latest version with major release number 1

Doing this generates a line in go.mod that specifies the version of the lib needed.

It also creates go.sum (checksum) which contains checksums representing the content of the required packages. Now I can run 'go mod verify' to verify the checksums of the downloaded packages on my machine match the entries in go.sum. It also helps with recreating dependency env.

CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
-- Important: Make sure to swap 'pass' with a password of your own choosing.
ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';

## Why use placeholder parameters instead of string interpolation

EX: 
```sql
INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
```
then execute the query with:

```go
result, err := m.DB.Exec(stmt, title, content, expires)
```

The DB.Exec() method avoids SQL injections by:
- creating a new prepared statement on the db using the provided sql query string.
- passing the parameter values to the db, the db then executes the prepared statement using these parameters. Since the params are transmitted later, after the statement has been compiled, the database treats them as pure data, so the intent of the statement can't change.
- it then deallocates the prepared statement on the database.

## row.Scan()
```go
err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
```

If any of the values are null in the db, it will raise an error since Go doesn't convert null to its corresponding type. The best thing to do is not allow null values in the database. Otherwise, you need to make the type of the object's attributes sql.NullType, for example sql.NullFloat, to handle that. 

# Go HTML Templating

When invoking a template from another template, you have the ability to specify which data object is passed.

So in the Go code, you 

```go
// files is list of file paths
ts, err := template.ParseFiles(files...)
// is an object with attrs and methods 
// that can be referenced in the template
err = ts.ExecuteTemplate(w, "base", data)
```

Then html 

```html
{{template "main" .}}
```

The dot is used to specify the entire object passed to the template
should be passed (pipelined) to the referenced template. You could specify an attribute to pass instead (e.g. .Name instead of .)

## Template actions and functions

- {{define "name"}} <h1>Name</h1>{{end}} defines a new html template
- {{template "name" .}} reference a defined html template for insertion (dot passes the data object passed to ExecuteTemplate explained above)
- {{block "name" .}} HTML {{end}} defines and uses the template
- {{if .Foo}} C1 {{else}} C2 {{end}} 
- {{with .Foo}} C1 {{else}} C2 {{end}} sets the dot to .Foo for the content of C1. If .Foo is empty, default to C2 (see view.tmpl for example)
- {{range .Foo}} C1 {{else}} C2 {{end}} iterate over .Foo, setting dot to each element then render C1. If empty render C2.

NOTE: {{else}} is optional in these cases
Also range loops can use break/continue
```html
{{range .Foo}}
    // End the loop if the .ID value equals 99.
    {{if eq .ID 99}}
        {{break}}
    {{end}}
    // ...
{{end}}
```


# Middleware

A Go web app is pretty much a chain of *ServeHTTP()* methods being called one after another

When our server receives a new HTTP request, it calls the servemux's ServeHTTP() method which looks up the relevant handler based on the request URL path, and then calls the handler's ServeHTTP() method. Middleware is just adding more handlers to this chain where it executes some logic (like logging a request or compressing a response) and then calls the ServeHTTP() method of the next handler in the chain.

The standard pattern for creating middleware:

```go

func myMiddleware(next http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        // TODO: Execute our middleware logic here...
        next.ServeHTTP(w, r)
    }
    return http.HandlerFunc(fn)
}
```

So myMiddleware is essentially a wrapper around the next handler (think decorator in Python). A closure is created over the next handler by creating an anonymous function. 

http.HandlerFunc is **NOT** a function but a type declaration of the type func(http.ResponseWriter, *http.Request) and we return an instance of this function type. You may notice the function signature doesn't return the http.HandlerFunc type, but that's because http.HandlerFunc fulfills the http.Handler interface by defining a ServeHTTP method with the correct parameters.

This can be abbreviated by:

```go
func myMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // TODO: Execute our middleware logic here...
        next.ServeHTTP(w, r)
    })
}
```

## Example: Setting security headers

```
Content-Security-Policy: default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com
Referrer-Policy: origin-when-cross-origin
X-Content-Type-Options: nosniff
X-Frame-Options: deny
X-XSS-Protection: 0
```

```go
func secureHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Note: This is split across multiple lines for readability. You don't
        // need to do this in your own code.
        w.Header().Set("Content-Security-Policy",
            "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
        w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "deny")
        w.Header().Set("X-XSS-Protection", "0")
        next.ServeHTTP(w, r)
    })
}
```

We want this to be executed on every request. So we want to position it before the servemux (the handler who's job it is to execute the correct handler depending on the route of the request).

secureHeaders -> servemux -> application handler

And then when the final function in the chain returns, you can have code execute at each handler on the way back up as seen below.

```go
func myMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Any code here will execute on the way down the chain.
        next.ServeHTTP(w, r)
        // Any code here will execute on the way back up the chain.
    })
}
```

Also if you call return before next.ServeHTTP(w,r), it ends the chain. For instance you could do an authorization check, if it fails, return before next and 