package mngment

import (
	"time"

	"github.com/google/uuid"
)

// Get fetch in the database the token associated with a specified user.
// `t`: is the token string.
// `u`:
func Get(t string, u string) {

}

type UserToken struct {
	Token    string
	UserID   string
	Expiracy time.Time
	IP       string
}

// Add adds a new token into the database.
func (ut *UserToken) Add() {
	ut.Token = uuid.New().String()
	ut.setExpiracy()
}

// Refresh update the token expiracy to be sure it doesn't expire.
func (ut *UserToken) Refresh() {
	ut.setExpiracy()

}

func (ut *UserToken) setExpiracy() {
	exp := time.Now()
	exp = exp.AddDate(0, 0, 14)
	ut.Expiracy = exp
}
