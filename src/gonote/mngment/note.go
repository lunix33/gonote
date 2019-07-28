package mngment

import (
	"gonote/db"
	"gonote/util"
	"html"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Note contains the content of a note.
type Note struct {
	ID      string
	Title   string
	UserID  string
	Public  bool
	Added   time.Time
	Deleted bool
}

// Delete trash or delete a note from the database.
// Struct required field: ID.
//
// "c" is an optional database connection.
//
// Returns any error (e) occured.
func (n *Note) Delete(c *db.Conn) (e error) {
	if n.ID == "" {
		e = errors.New("missing field to be able to delete the note")
		return
	}

	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{n.ID}
		if !n.Deleted {
			_, _, e = db.Run(c, noteTrashQuery, p, nil)
		} else {
			_, _, e = db.Run(c, noteDeleteQuery, p, nil)
		}
	})

	return e
}

// Add adds a note into the database.
// Struct required field: Title, UserID, Public.
//
// "c" is an optional database connection.
//
// Returns any error (e) occured.
func (n *Note) Add(c *db.Conn) (e error) {
	if n.UserID == "" {
		e = errors.New("missing field to be able to add the note")
		return
	}

	// Generate Object
	n.ID = uuid.New().String()
	n.Title = html.EscapeString(n.Title)
	n.Added = time.Now()
	n.Deleted = false

	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{n.ID, n.Title, n.UserID, n.Public}
		_, _, e = db.Run(c, noteAddQuery, p, nil)
		if e != nil {
			e = errors.Wrap(e, "unable to add the note")
		}
	})

	return e
}

// Update change the content and save to file.
//
// "n2" is the update document.
// "c" is an optional database connection.
//
// Returns any error (e) occurred.
func (n *Note) Update(n2 *Note, c *db.Conn) (e error) {
	if n.ID == "" {
		e = errors.New("missing field to be able to update the note")
		return
	}

	// Update the object.
	n.Title = html.EscapeString(n2.Title)
	n.Public = n2.Public

	db.MustConnect(c, func(id *db.Conn) {
		p := []interface{}{n.Title, n.Public, n.ID}
		_, _, e = db.Run(id, noteUpdateQuery, p, nil)
		if e != nil {
			e = errors.Wrap(e, "unable to update note")
		}
	})

	return e
}

// GetTags gets the tags associated with the note.
//
// "c" is an optional database connection.
//
// Returns the list of tags (t) associated with the note.
func (n *Note) GetTags(c *db.Conn) (t []*Tag) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{n.ID}
		rst, _, err := db.Run(c, noteGetAllTagsQuery, p, reflect.TypeOf(Tag{}))
		if err == nil {
			for _, v := range rst {
				t = append(t, v.(*Tag))
			}
		} else {
			util.LogErr(err)
		}
	})

	return t
}

// GetNoteContent gets all the contents of the note.
//
// "c" is an optional database connection.
//
// Returns a list of note content (nc) associated with the note.
func (n *Note) GetNoteContent(c *db.Conn) (nc []*NoteContent) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{n.ID}
		rst, _, err := db.Run(c, noteGetNoteContentQuery, p, reflect.TypeOf(NoteContent{}))
		if err == nil {
			for _, v := range rst {
				nc = append(nc, v.(*NoteContent))
			}
		} else {
			util.LogErr(err)
		}
	})

	return nc
}
