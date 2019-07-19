package mngment

import (
	"gonote/db"
	"reflect"
	"time"
)

// GetNoteContent fetch a specific version of a note content.
// `id` is the ID of a note.
// `vers` is the version of the note content.
// `c` is an optional database connection.
// Returns the content of a note.
func GetNoteContent(id string, vers string, c *db.Conn) (nc *NoteContent) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{id, vers}
		rst, cnt, err := db.Run(c, noteContentGetQuery, p, reflect.TypeOf(NoteContent{}))
		if err == nil && cnt > 0 {
			nc = rst[0].(*NoteContent)
		}
	})

	return nc
}

// GetNoteContents gets all the version content of a specific note.
// `id` is the id of a note.
// `c` is an optional database connection.
// Returns a list of all the version of a note.
func GetNoteContents(id string, c *db.Conn) (ncs []*NoteContent) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{id}
		res, _, err := db.Run(c, noteContentGetAllQuery, p, reflect.TypeOf(NoteContent{}))
		if err == nil {
			for _, v := range res {
				ncs = append(ncs, v.(*NoteContent))
			}
		}
	})

	return ncs
}

// TODO SEARCH

// NoteContent represent the content of a Note.
type NoteContent struct {
	NoteID  string
	Version int
	Content string
	Updated time.Time
}

// Add adds a new version of a note content.
// `c` is an optional database connection.
// Returns any errors (e) occurred.
func (nc *NoteContent) Add(c *db.Conn) (e error) {
	db.MustConnect(c, func(c *db.Conn) {
		// Update the struct fields.
		nc.setLastVersion(c)
		nc.Updated = time.Now()

		// Run the insert.
		p := []interface{}{nc.NoteID, nc.Version, nc.Content, nc.Updated}
		_, _, e = db.Run(c, noteContentInsertQuery, p, nil)
	})

	return e
}

// TODO DELETE

// TODO UPDATE

// setLastVersion set the version of the note as the last version.
// `c` is a database connection.
func (nc *NoteContent) setLastVersion(c *db.Conn) {
	nts := GetNoteContents(nc.NoteID, c)
	if len(nts) > 0 {
		nc.Version = nts[0].Version + 1
	} else {
		nc.Version = 1
	}
}
