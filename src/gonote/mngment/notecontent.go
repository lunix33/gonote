package mngment

import (
	"gonote/db"
	"html"
	"time"
)

// NoteContent represent the content of a Note.
type NoteContent struct {
	NoteID  string
	Version int
	Content string
	Updated time.Time
}

// Add adds a new version of a note content.
// Struct required fields: NoteID, Content.
// `c` is an optional database connection.
// Returns any errors (e) occurred.
func (nc *NoteContent) Add(c *db.Conn) (e error) {
	db.MustConnect(c, func(c *db.Conn) {
		// Update the struct fields.
		nc.Version = (nc.getLastVersionNumber(c) + 1)
		nc.Updated = time.Now()
		nc.Content = html.EscapeString(nc.Content)

		// Run the insert.
		p := []interface{}{nc.NoteID, nc.Version, nc.Content, nc.Updated}
		_, _, e = db.Run(c, noteContentInsertQuery, p, nil)
	})

	return e
}

// Delete removes a specified version of a note content from the database.
// Struct required fields: NoteID, Version.
// `c` is an optional database connection.
// Returns any errors (e) occured.
func (nc *NoteContent) Delete(c *db.Conn) (e error) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{nc.NoteID, nc.Version}
		_, _, e = db.Run(c, noteContentDeleteQuery, p, nil)
	})

	return e
}

// Update update the latest version of a note content.
// If the user is trying to edit a old version, a new version will be created.
// `c` is an optional database configuration.
// Returns any errors (e) occured.
func (nc *NoteContent) Update(c *db.Conn) (e error) {
	db.MustConnect(c, func(c *db.Conn) {
		// User can only update the latest revision.
		// If they are making an edit on a old revision, then a new revision is created.
		lstVers := nc.getLastVersionNumber(c)
		if lstVers == nc.Version {
			nc.Updated = time.Now()
			nc.Content = html.EscapeString(nc.Content)
			p := []interface{}{nc.Content, nc.Updated, nc.NoteID, nc.Version}
			_, _, e = db.Run(c, noteContentUpdateQuery, p, nil)
		} else {
			e = nc.Add(c)
		}
	})

	return e
}

// getLastVersionNumber gets the last version number of the note.
// `c` is a database connection.
// Returns the last version number.
func (nc *NoteContent) getLastVersionNumber(c *db.Conn) int {
	nts := GetAllNoteContents(nc.NoteID, c)
	if len(nts) > 0 {
		return nts[0].Version
	}
	return 1
}
