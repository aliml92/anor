<img src="./assets/logo.svg" alt="anor" width="100" height="100">

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/aliml92/anor/blob/master/LICENSE)

**Anor** is a fullstack ecommerce app built with *go/htmx* (in heavy development process).

* The project uses:
    * [go](https://github.com/golang/go) for backend
    * [htmx](https://github.com/bigskysoftware/htmx) and [_hyperscript](https://github.com/bigskysoftware/_hyperscript) for fronted
    * [pgx](https://github.com/jackc/pgx) as the database driver
    * [sqlc](https://github.com/kyleconroy/sqlc) to generate Go code from SQL queries
    * [sendinblue](https://github.com/sendinblue/APIv3-go-library) for transactional emails
    * [go-typesense](https://github.com/aliml92/go-typesense) as Typesense client
    * [goose](https://github.com/pressly/goose) to manage database migrations
    * [task](https://github.com/go-task/task) as a task runner 

# Getting started 
Make sure you have *docker/docker compose*, *sqlc*, *goose* and *task* installed in your machine

First start *Postgres* and *Typesense*:
```
    task compose-up
``` 
Then, run database migrations:
```
    task goose-up
```
and import sample data into database using:
```
    task import-dataset
```

Now you can start the app:
```
    export CONFIG_FILEPATH=./config.dev.yaml
    go run cmd/anor/*.go
```
Project starts on port 8008 by default.

# Testing
Coming soon