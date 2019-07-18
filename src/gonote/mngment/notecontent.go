package mngment

import "time"

// NoteContent represent the content of a Note.
type NoteContent struct {
	NoteID  string
	Version int
	Content string
	Updated time.Time
}
