package mngment

import (
	"gonote/db"
	"reflect"

	"github.com/google/uuid"
)

// GetUser get a user from the database.
// `uname` is the username of the user.
// `c` is an optional database connection.
// Returns the user (u) found. Will be nil if an error occure or no user is found.
func GetUser(uname string, c *db.Conn) (u *User) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{uname}
		rst, cnt, err := db.Run(c, userGetQuery, p, reflect.TypeOf(User{}))
		if err == nil && cnt > 0 {
			u = rst[0].(*User)
		}
	})

	return u
}

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
// `c` is an optional database connection.
// Returns any error (e) occured.
func (u *User) Add(c *db.Conn) (e error) {
	// If the password was set, encode it.
	if u.Password != "" {
		u.SetPassword(u.Password)
	}

	u.ID = uuid.New().String()
	u.Deleted = false

	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{u.ID, u.Username, u.Password, u.Email, u.IsAdmin}
		_, _, e = db.Run(c, userAddQuery, p, nil)
	})

	return e
}

// Update updates the user with data of a new user object.
// `c` is an optional database connection.
// `uu` is the user with the update information
// Returns any error (e) occured.
func (u *User) Update(uu *User, c *db.Conn) (e error) {
	// Update password if not empty
	if uu.Password != "" {
		u.SetPassword(uu.Password)
	}

	// Update email if not empty
	if uu.Email != "" {
		u.Email = uu.Email
	}

	if uu.Username != "" {
		u.Username = uu.Username
	}

	u.IsAdmin = uu.IsAdmin

	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{u.Username, u.Password, u.Email, u.IsAdmin, u.ID}
		_, _, e = db.Run(c, userUpdateQuery, p, nil)
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
	})

	return e
}
