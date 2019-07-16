package tag

import (
	"gonote/db"
)

// Tag represent a Note category tag.
type Tag struct {
	NoteID string
	Name   string
}

// Add adds a new tag.
// `dbID` is the database ID.
// Returns any error occured.
func (t *Tag) Add(dbID *string) (e error) {
	db.MustConnect(dbID, func(id string) {
		q := `INSERT INTO Tag (Name) VALUES (?)`
		p := []interface{}{t.Name}
		_, _, qErr := db.Run(id, q, p, nil)
		if qErr != nil {
			e = qErr
		}
	})

	return e
}

// GetNotes gets all the notes associated with a tag name.
func (t *Tag) GetNotes(dbID *string) (r []*Note) {
	db.MustConnect(dbID, func(id string) {
		q := `
			SELECT Note.* FROM Note
			INNER JOIN NoteTag ON Note.ID = NoteTag.NoteID
			WHERE NoteTag.TagName = ?`
		p := []interface{}{t.Name}
	})
	return
}
