package mngment

import (
	"gonote/db"
	"html"
	"time"

	"github.com/pkg/errors"
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
//
// "c" is an optional database connection.
//
// Returns any errors (e) occurred.
func (nc *NoteContent) Add(c *db.Conn) (e error) {
	if nc.NoteID == "" {
		errors.New("missing field to be able to add the note content")
		return
	}

	db.MustConnect(c, func(c *db.Conn) {
		// Update the struct fields.
		nc.Version = (nc.getLastVersionNumber(c) + 1)
		nc.Updated = time.Now()
		nc.Content = html.EscapeString(nc.Content)

		// Run the insert.
		p := []interface{}{nc.NoteID, nc.Version, nc.Content, nc.Updated}
		_, _, e = db.Run(c, noteContentInsertQuery, p, nil)
		if e != nil {
			e = errors.Wrap(e, "was unable to add the note content")
		}
	})

	return e
}

// Delete removes a specified version of a note content from the database.
// Struct required fields: NoteID, Version.
//
// "c" is an optional database connection.
//
// Returns any errors (e) occured.
func (nc *NoteContent) Delete(c *db.Conn) (e error) {
	if nc.NoteID == "" || nc.Version == 0 {
		e = errors.New("missing field to be able to delete the note content")
		return
	}
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{nc.NoteID, nc.Version}
		_, _, e = db.Run(c, noteContentDeleteQuery, p, nil)
		if e != nil {
			e = errors.Wrap(e, "was unable to delete the note content")
		}
	})

	return e
}

// Update update the latest version of a note content.
// If the user is trying to edit a old version, a new version will be created.
//
// "nc2" is the update document.
// "c" is an optional database configuration.
//
// Returns any errors (e) occured.
func (nc *NoteContent) Update(nc2 *NoteContent, c *db.Conn) (e error) {
	if nc.NoteID == "" || nc.Version == 0 {
		e = errors.New("missing field to be able to update the note content")
		return
	}

	db.MustConnect(c, func(c *db.Conn) {
		// User can only update the latest revision.
		// If they are making an edit on a old revision, then a new revision is created.
		lstVers := nc.getLastVersionNumber(c)
		if lstVers == nc.Version {
			nc.Updated = time.Now()
			nc.Content = html.EscapeString(nc2.Content)
			p := []interface{}{nc.Content, nc.Updated, nc.NoteID, nc.Version}
			_, _, e = db.Run(c, noteContentUpdateQuery, p, nil)
		} else {
			e = nc.Add(c)
		}

		if e != nil {
			e = errors.Wrap(e, "unable to update the note content.")
		}
	})

	return e
}

// getLastVersionNumber gets the last version number of the note.
//
// "c" is a database connection.
//
// Returns the last version number.
func (nc *NoteContent) getLastVersionNumber(c *db.Conn) int {
	nts := GetAllNoteContents(nc.NoteID, c)
	if len(nts) > 0 {
		return nts[0].Version
	}
	return 1
}
