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

