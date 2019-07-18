package mngment

import (
	"gonote/db"
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
		}
	})

	return t
}

// Tag represent a Note category tag.
type Tag struct {
	NoteID string
	Name   string
}

// Add adds a new tag.
// `dbID` is the database ID.
// Returns any error occured.
func (t *Tag) Add(c *db.Conn) (e error) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{t.NoteID, t.Name}
		_, _, e = db.Run(c, tagAddQuery, p, nil)
	})

	return e
}

// Remove removes the tag association from the database.
// `c` is an optional database connection.
// Returns any error (e) occured.
func (t *Tag) Remove(c *db.Conn) (e error) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{t.NoteID, t.Name}
		_, _, e = db.Run(c, tagRemoveQuery, p, nil)
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
