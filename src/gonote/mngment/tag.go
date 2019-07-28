package mngment

import (
	"gonote/db"

	"github.com/pkg/errors"
)

// Tag represent a Note category tag.
type Tag struct {
	NoteID string
	Name   string
}

// Add adds a new tag.
//
// "c" is the database ID.
//
// Returns any error occured.
func (t *Tag) Add(c *db.Conn) (e error) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{t.NoteID, t.Name}
		_, _, e = db.Run(c, tagAddQuery, p, nil)
		if e != nil {
			e = errors.Wrapf(e, "unable to add tag \"%s\" to note %s", t.Name, t.NoteID)
		}
	})

	return e
}

// Remove removes the tag association from the database.
//
// "c" is an optional database connection.
//
// Returns any error (e) occured.
func (t *Tag) Remove(c *db.Conn) (e error) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{t.NoteID, t.Name}
		_, _, e = db.Run(c, tagRemoveQuery, p, nil)
		if e != nil {
			e = errors.Wrapf(e, "unable to remove tag \"%s\" from note %s", t.Name, t.NoteID)
		}
	})

	return e
}

// TODO: Needs to find the actual query first.
// GetNotes gets all the notes associated with a tag name.
// func (t *Tag) GetNotes(c *db.Conn) (n []*Note) {
// 	db.MustConnect(c, func(c *db.Conn) {
// 		p := []interface{}{t.NoteID}
// 		rst, _, err := db.Run(c, tagGetNotesQuery, p, reflect.TypeOf(Note{}))
// 	})
// 	return n
// }
