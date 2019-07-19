package mngment

import (
	"gonote/db"
	"reflect"
)

// GetNote retrive a note with a specified id.
// If the note ain't found, returns nil.
func GetNote(nID string, c *db.Conn) (n *Note) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{nID}
		rst, cnt, err := db.Run(c, noteGetQuery, p, reflect.TypeOf(Note{}))
		if err == nil && cnt > 0 {
			n = rst[0].(*Note)
		}
	})

	return n
}

// SearchNotes fetch all the notes from the DB which correspond to criterias.
func SearchNotes(author string, deleted bool, public bool, search string, c *db.Conn) (sr []*Note) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{author, public, deleted, search, search}
		rst, _, err := db.Run(c, noteSearchQuery, p, reflect.TypeOf(Note{}))
		if err == nil {
			for _, v := range rst {
				sr = append(sr, v.(*Note))
			}
		}
	})

	return sr
}
