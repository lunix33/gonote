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

// TokenCleanupRoutine is a long running routine made to cleanup the invalid tokens every 12 hours.
func TokenCleanupRoutine() {
	for {
		log.Println("Runinng Token cleanup...")

		var cleanupCount int
		db.MustConnect(nil, func(c *db.Conn) {
			// TODO: Get all the user tokens, and validate every token.
			//       Count the number of token deleted with "cleanupCount".
			//       Run every validation with "dry" to true.
		})

		log.Printf("%d tokens were cleanned up...\n", cleanupCount)

		// Execute exery 12 hours.
		time.Sleep(12 * time.Hour)
	}
}
