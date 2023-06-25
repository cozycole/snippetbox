# Creating a module
Module path: a standard name or identifier for the project. Almost any string can be used but it needs to be unique to avoid import conflicts with other projects. Common convention is to base a module path on a URL you own. Example: snippetbox.cozycole.net

You then initialize the go module with that name 

```bash
go mod init snippetbox.cozycole.net
```

This generates a go.mod file which signals the directory it is stored in is a go module.

