package user

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Bloggers gets all the possible people who can blog on our site.
func Bloggers(ctx context.Context, query string, db *sqlx.DB) ([]User, error) {
	usrs := []User{}
	qContext, _ := context.WithDeadline(ctx, time.Now().Add(time.Duration(2)*time.Second))
	err := db.SelectContext(qContext, &usrs, "SELECT * FROM bloggers")

	if err != nil {
		return []User{}, err
	}

	return usrs, err

}

// SpecificBlogger attempts to get a single user who can blog on our site via their username.
// We rely on the schema to enforce uniqueness
func SpecificBlogger(ctx context.Context, identifier string, db *sqlx.DB) ([]User, error) {
	usrs := []User{}
	qContext, _ := context.WithDeadline(ctx, time.Now().Add(time.Duration(2)*time.Second))
	err := db.SelectContext(qContext, &usrs, fmt.Sprintf("SELECT * FROM bloggers WHERE username = '%s'", identifier))

	if err != nil {
		return []User{}, err
	}

	return usrs, err

}

const uInsert = "INSERT INTO `bloggers`( `username`, `pass_hash`, `first_name`, `last_name`) VALUES(':username', ':pass_hash', ':first_name', ':last_name');"

// InsertBlogger into the data store. Return an error if for some reason it is invalid.
// This perhaps should be a lot more verbose of an error with some extra user facing intent.
func InsertBlogger(ctx context.Context, user interface{}, db *sqlx.DB) error {
	fmt.Println(uInsert, fmt.Sprintf("%v", user))
	// qContext, _ := context.WithDeadline(ctx, time.Now().Add(time.Duration(2)*time.Second))
	_, err := db.NamedExec(uInsert, user)
	return err
}
