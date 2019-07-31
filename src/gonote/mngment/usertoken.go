package mngment

import (
	"gonote/db"
	"gonote/util"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	// LoginToken is a type of token made for user login.
	LoginToken = "login"

	// PasswordResetToken is a type of token made for user's password reset.
	PasswordResetToken = "passreset"
)

// UserToken represents the login token of a user.
type UserToken struct {
	Token   string
	Type    string
	UserID  string
	Created time.Time
	Expiry  time.Time
	IP      string
}

// Add adds a new token into the database.
//
// "c" is an optional database connection.
//
// Returns any error (e) occured.
func (ut *UserToken) Add(c *db.Conn) (e error) {
	// Set default values of a Token.
	if ut.Type == "" {
		ut.Type = LoginToken
	}
	ut.Token = uuid.New().String()
	ut.Created = time.Now()
	ut.setExpiry()

	// IP must be set.
	if ut.IP == "" {
		return errors.New("the IP of the client must be set")
	}

	// Insert the token in the DB.
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{ut.Token, ut.Type, ut.UserID, ut.Expiry, ut.IP}
		_, _, e = db.Run(c, userTokenInsertQuery, p, nil)
		if e != nil {
			e = errors.Wrapf(e, "unable to create a new token for user %s", ut.UserID)
		}
	})

	return e
}

// Refresh update the token Expiry to be sure it doesn't expire.
//
// "c" is an optional database connection.
//
// Returns any error (e) occured.
func (ut *UserToken) Refresh(c *db.Conn) (e error) {
	// PasswordResetToken can't be refreshed.
	if ut.Type == PasswordResetToken {
		return nil
	}

	ut.setExpiry()

	// IP must be set.
	if ut.IP == "" {
		return errors.New("the IP of the client must be set")
	}

	// Update the token with a new IP and a new expiration.
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{ut.Expiry, ut.IP, ut.UserID, ut.Token}
		_, _, e = db.Run(c, userTokenRefreshQuery, p, nil)
		if e != nil {
			e = errors.Wrapf(e, "unable to refresh token %s", ut.Token)
		}
	})

	return e
}

// Delete remove the token from the database.
//
// "c" is an optional database connection.
//
// Returns any error (e) occured.
func (ut *UserToken) Delete(c *db.Conn) (e error) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{ut.UserID, ut.Token}
		_, _, e = db.Run(c, userTokenDeleteQuery, p, nil)
		if e != nil {
			e = errors.Wrapf(e, "unable to delete the token %s", ut.Token)
		}
	})
	return e
}

// Validate verify if a token is still valid.
// It also delete invalid tokens and one-time tokens (when used)
//
// "dry" True if the application should dry run the validation.
// "c" is an optional database connection.
//
// Returns wether or not a token is valid (v).
func (ut *UserToken) Validate(dry bool, c *db.Conn) (v bool) {
	var (
		now = time.Now()
		err error
	)
	if now.Before(ut.Expiry) {
		v = true

		// Remove One time token
		if ut.Type == PasswordResetToken && !dry {
			err = ut.Delete(c)
		}
	} else {
		// Token is expired.
		v = false
		if !dry {
			err = ut.Delete(c)
		}
	}

	if err != nil {
		util.LogErr(errors.Wrapf(err, "an error occurred while trying to remove an invalid token (%s)", ut.Token))
	}

	return v
}

// setExpiry sets the token's Expiry based on the token type.
func (ut *UserToken) setExpiry() {
	exp := time.Now()

	// Set Expiry based on token type.
	if ut.Type == PasswordResetToken {
		// One day later.
		exp = exp.AddDate(0, 0, 1)
	} else {
		// 2 weeks later.
		exp = exp.AddDate(0, 0, 14)
	}

	ut.Expiry = exp
}
