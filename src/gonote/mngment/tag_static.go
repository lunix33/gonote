package mngment

import (
	"gonote/db"
	"gonote/util"
	"reflect"
)

// GetAllTags gets a list of all the available tags.
// `c` is an optional database connection.
// Returns a list of tags (t).
func GetAllTags(c *db.Conn) (t []*Tag) {
	db.MustConnect(c, func(c *db.Conn) {
		tags, _, err := db.Run(c, tagGetAllQuery, nil, reflect.TypeOf(Tag{}))
		if err == nil {
			for _, v := range tags {
				t = append(t, v.(*Tag))
			}
		} else {
			util.LogErr(err)
		}
	})

	return t
}
