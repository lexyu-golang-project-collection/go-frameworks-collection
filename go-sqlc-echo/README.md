```
.
|-- Makefile
|-- bin
|   `-- api
|-- cmd
|   `-- main.go
|-- config
|   |-- constants.go
|   |-- database.go
|   `-- env.go
|-- go.mod
|-- go.sum
|-- internal
|   |-- db
|   |   |-- embed
|   |   |   |-- sql.go
|   |   |   `-- sqlc
|   |   |       |-- mysql
|   |   |       |   |-- query.sql
|   |   |       |   `-- schema.sql
|   |   |       |-- postgres
|   |   |       |   |-- query.sql
|   |   |       |   `-- schema.sql
|   |   |       `-- sqlite
|   |   |           |-- query.sql
|   |   |           `-- schema.sql
|   |   |-- mysql
|   |   |   |-- db.go
|   |   |   |-- models.go
|   |   |   `-- query.sql.go
|   |   |-- postgres
|   |   |   |-- db.go
|   |   |   |-- models.go
|   |   |   `-- query.sql.go
|   |   `-- sqlite
|   |       |-- db.go
|   |       |-- models.go
|   |       `-- query.sql.go
|   |-- handler
|   |   |-- author_handler.go
|   |   |-- base.go
|   |   |-- handler.go
|   |   `-- task_handler.go
|   |-- model
|   |   `-- author.go
|   |-- repository
|   |   `-- manager.go
|   `-- service
|       |-- author_service.go
|       |-- service.go
|       `-- tracker.go
|-- pkg
|   `-- middleware
|       `-- server_state.go
|-- router
|   `-- router.go
`-- sqlc.yaml
```