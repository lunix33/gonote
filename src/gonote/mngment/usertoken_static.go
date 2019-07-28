package mngment

import (
	"gonote/db"
	"gonote/util"
	"reflect"
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
