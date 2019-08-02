package mngment

import (
	"gonote/db"
	"gonote/util"
	"log"
	"reflect"
	"time"
)

// GetUserToken fetch in the database the token associated with a specified user.
//
// "t" is the token string.
// "u" is the user id of the user owning the token.
// "c" is an optional database connection.
//
// Returns the user token (r) respecting the constraints.
func GetUserToken(t string, u string, c *db.Conn) (r *UserToken) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{u, t}
		rst, cnt, err := db.Run(c, userTokenGetQuery, p, reflect.TypeOf(UserToken{}))
		if err == nil && cnt > 0 {
			r = rst[0].(*UserToken)
		} else if err != nil {
			util.LogErr(err)
		}
	})

	return r
}

// GetAllUserTokens gets all the users tokens.
//
// "c" is an optional database connection.
//
// Returns a list of user tokents (t).
func GetAllUserTokens(c *db.Conn) (t []*UserToken) {
	db.MustConnect(c, func(c *db.Conn) {
		rst, _, err := db.Run(c, userTokenGetAllQuery, nil, reflect.TypeOf(UserToken{}))
		if err == nil {
			for _, v := range rst {
				t = append(t, v.(*UserToken))
			}
		}
	})

	return t
}

// TokenCleanupRoutine is a long running routine made to cleanup the invalid tokens every 12 hours.
func TokenCleanupRoutine() {
	for {
		log.Println("Runinng Token cleanup...")

		var cleanupCount int
		db.MustConnect(nil, func(c *db.Conn) {
			uts := GetAllUserTokens(c)
			for _, ut := range uts {
				if !ut.Validate(true, c) {
					err := ut.Delete(c)
					if err == nil {
						cleanupCount++
					}
				}
			}
		})

		log.Printf("... %d tokens were cleanned up.\n", cleanupCount)

		// Execute exery 12 hours.
		time.Sleep(12 * time.Hour)
	}
}
