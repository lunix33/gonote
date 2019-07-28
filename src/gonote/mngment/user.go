package mngment

import (
	"gonote/db"
	"html"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// User represent a user of the platform.
type User struct {
	ID       string
	Username string
	Password string
	Email    string
	Deleted  bool
	IsAdmin  bool
}

// Add adds a user.
//
// "c" is an optional database connection.
//
// Returns any error (e) occured.
func (u *User) Add(c *db.Conn) (e error) {
	u.ID = uuid.New().String()
	u.Username = html.EscapeString(u.Username)
	u.SetPassword(u.Password)
	u.Email = html.EscapeString(u.Email)
	u.Deleted = false

	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{u.ID, u.Username, u.Password, u.Email, u.IsAdmin}
		_, _, e = db.Run(c, userAddQuery, p, nil)
		if e != nil {
			e = errors.Wrap(e, "unable to create the new user")
		}
	})

	return e
}

// Update updates the user with data of a new user object.
//
// "uu" is the user with the update information
// "c" is an optional database connection.
//
// Returns any error (e) occured.
func (u *User) Update(uu *User, c *db.Conn) (e error) {
	if uu.Username != "" {
		u.Username = html.EscapeString(uu.Username)
	}

	// Update password if not empty
	if uu.Password != "" {
		u.SetPassword(uu.Password)
	}

	// Update email if not empty
	if uu.Email != "" {
		u.Email = html.EscapeString(uu.Email)
	}

	u.IsAdmin = uu.IsAdmin

	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{u.Username, u.Password, u.Email, u.IsAdmin, u.ID}
		_, _, e = db.Run(c, userUpdateQuery, p, nil)
		if e != nil {
			e = errors.Wrapf(e, "unable to update the user %s (%s)", u.Username, u.ID)
		}
	})

	return e
}

// Delete makes the user inactive
// `c` is an optional database connection.
// Returns any error (e) occured.
func (u *User) Delete(c *db.Conn) (e error) {
	u.Deleted = true
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{u.ID}
		_, _, e = db.Run(c, userDeleteQuery, p, nil)
		if e != nil {
			e = errors.Wrapf(e, "unable to delete the user %s (%s)", u.Username, u.ID)
		}
	})

	return e
}
