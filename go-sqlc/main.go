package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"reflect"

	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/tutorial"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

func run() error {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	queries := tutorial.New(db)

	// list all authors
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	// create an author
	insertedAuthor, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		return err
	}
	log.Println(insertedAuthor)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthor.ID)
	if err != nil {
		return err
	}
	// prints true
	log.Println(reflect.DeepEqual(insertedAuthor, fetchedAuthor))

	return err
}

func main() {

	if err := run(); err != nil {
		log.Fatal(err)
	}

}
