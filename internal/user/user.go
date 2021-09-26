package user

import (
	"context"

	"github.com/implausiblyfun/blogforge/internal/encrypt"
	"github.com/jmoiron/sqlx"
)

// User captures what it means to be a blogger in our blogspace.
type User struct {
	ID         int    `json:"id" db:"id"`                 //`id` int(11) NOT NULL AUTO_INCREMENT,
	Username   string `json:"username" db:"username"`     //`username` VARCHAR(255) NOT NULL,
	FirstName  string `json:"first_name" db:"first_name"` //`first_name` VARCHAR(255),
	LastName   string `json:"last_name" db:"last_name"`   //`last_name`  VARCHAR(255),
	HashedPass string `json:"pass_hash" db:"pass_hash"`   //`pass_hash` VARCHAR(255),
	Password   string `json:"password" db:"-"`            //`pass_hash` VARCHAR(255),
}

// ValidatePassword for a user.
func (u *User) ValidatePassword(pass string) bool {
	return u.HashedPass == encrypt.Encode(pass)
}

// NewUser from the given options.
func NewUser(username, first, last, pass string) *User {
	return &User{
		Username:  username,
		FirstName: first,
		LastName:  last,
		Password:  pass,
	}
}

// Store the user to the database.
func (u *User) Store(ctx context.Context, db *sqlx.DB, overwrite bool) error {
	if u.HashedPass == "" {
		u.HashedPass = encrypt.Encode(u.Password)
	}
	var err error
	if overwrite {
		err = InsertBlogger(ctx, u, db)
	}
	return err
}
