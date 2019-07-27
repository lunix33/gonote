package mngment

import (
	"gonote/db"
	"reflect"
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

// GetUserByID get a user from the database by the id.
// `id` is the user id.
// `c` is an optional database connection.
// Returns the user (u) found. Will be nil if an error occure or no user is found.
func GetUserByID(id string, c *db.Conn) (u *User) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{id}
		rst, cnt, err := db.Run(c, userGetByIDQuery, p, reflect.TypeOf(User{}))
		if err == nil && cnt > 0 {
			u = rst[0].(*User)
		}
	})

	return u
}
