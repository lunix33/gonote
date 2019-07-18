package mngment

import "gonote/db"

// GetUser get a user from the database.
func GetUser(c *db.Conn) (u *User) {}

// User represent a user of the platform.
type User struct {
	Username string
	Password string
	Email    string
	Deleted  bool
	IsAdmin  bool
}

// Add adds a user.
func (u *User) Add(c *db.Conn) (e error) {}

func (u *User) Update(c *db.Conn) (e error) {}

func (u *User) Delete(c *db.Conn) (e error) {}
