package mngment

import (
	"errors"
	"gonote/db"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// GetUserToken fetch in the database the token associated with a specified user.
//
// `t` is the token string.
// `u` is the user id of the user owning the token.
// `c` is an optional database connection.
//
// Returns the user token (r) respecting the constraints.
func GetUserToken(t string, u string, c *db.Conn) (r *UserToken) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{u, t}
		rst, cnt, err := db.Run(c, userTokenGetQuery, p, reflect.TypeOf(UserToken{}))
		if err == nil && cnt > 0 {
			r = rst[0].(*UserToken)
		}
	})

	return r
}

const (
	// LoginToken is a type of token made for user login.
	LoginToken = "login"

	// PasswordResetToken is a type of token made for user's password reset.
	PasswordResetToken = "passreset"
)

// UserToken represents the login token of a user.
type UserToken struct {
	Token    string
	Type     string
	UserID   string
	Created  time.Time
	Expiracy time.Time
	IP       string
}

// Add adds a new token into the database.
// `c` is an optional database connection.
// Returns any error (e) occured.
func (ut *UserToken) Add(c *db.Conn) (e error) {
	// Set default values of a Token.
	if ut.Type == "" {
		ut.Type = LoginToken
	}
	ut.Token = uuid.New().String()
	ut.Created = time.Now()
	ut.setExpiracy()

	// IP must be set.
	if ut.IP == "" {
		return errors.New("the IP of the client must be set")
	}

	// Insert the token in the DB.
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{ut.Token, ut.Type, ut.UserID, ut.Expiracy, ut.IP}
		_, _, e = db.Run(c, userTokenInsertQuery, p, nil)
	})

	return e
}

// Refresh update the token expiracy to be sure it doesn't expire.
// `c` is an optional database connection.
// Returns any error (e) occured.
func (ut *UserToken) Refresh(c *db.Conn) (e error) {
	ut.setExpiracy()

	// IP must be set.
	if ut.IP == "" {
		return errors.New("the IP of the client must be set")
	}

	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{ut.Expiracy, ut.IP, ut.UserID, ut.Token}
		_, _, e = db.Run(c, userTokenRefreshQuery, p, nil)
	})

	return e
}

// Delete remove the token from the database.
// `c` is an optional database connection.
// Returns any error (e) occured.
func (ut *UserToken) Delete(c *db.Conn) (e error) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{ut.UserID, ut.Token}
		_, _, e = db.Run(c, userTokenDeleteQuery, p, nil)
	})
	return e
}

// Validate verify if a token is still valid.
// It also delete invalid tokens and one-time tokens (when used)
// `c` is an optional database connection
// Returns wether or not a token is valid (v).
func (ut *UserToken) Validate(c *db.Conn) (v bool) {
	now := time.Now()
	if now.Before(ut.Expiracy) {
		v = true

		// Remove One time token
		if ut.Type == PasswordResetToken {
			ut.Delete(c)
		}
	} else {
		// Token is expired.
		v = false
		ut.Delete(c)
	}

	return v
}

// setExpiracy sets the token's expiracy based on the token type.
func (ut *UserToken) setExpiracy() {
	exp := time.Now()

	// Set expiracy based on token type.
	if ut.Type == PasswordResetToken {
		// One day later.
		exp = exp.AddDate(0, 0, 1)
	} else {
		// 2 weeks later.
		exp = exp.AddDate(0, 0, 14)
	}

	ut.Expiracy = exp
}
