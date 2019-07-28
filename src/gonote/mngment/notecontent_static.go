package mngment

import (
	"gonote/db"
	"log"
	"reflect"
)

// GetNoteContent fetch a specific version of a note content.
//
// "id" is the ID of a note.
// "vers" is the version of the note content.
// "c" is an optional database connection.
//
// Returns the content (nc) of a note or nil.
func GetNoteContent(id string, vers string, c *db.Conn) (nc *NoteContent) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{id, vers}
		rst, cnt, err := db.Run(c, noteContentGetQuery, p, reflect.TypeOf(NoteContent{}))
		if err == nil && cnt > 0 {
			nc = rst[0].(*NoteContent)
		} else if err != nil {
			log.Fatalln(err)
		}
	})

	return nc
}

// GetAllNoteContents gets all the version content of a specific note.
//
// "id" is the id of a note.
// "c" is an optional database connection.
//
// Returns a list of all the version of a note (ncs).
func GetAllNoteContents(id string, c *db.Conn) (ncs []*NoteContent) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{id}
		res, _, err := db.Run(c, noteContentGetAllQuery, p, reflect.TypeOf(NoteContent{}))
		if err == nil {
			for _, v := range res {
				ncs = append(ncs, v.(*NoteContent))
			}
		} else {
			log.Fatalln(err)
		}
	})

	return ncs
}
