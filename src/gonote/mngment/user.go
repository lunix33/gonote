package mngment

// GetUser get a user from the database.
func GetUser(id string) User {}

// User represent a user of the platform.
type User struct {
	Username string
	Password string
	Email    string
	Deleted  bool
}

// Add adds a user.
func (u *User) Add() {}
