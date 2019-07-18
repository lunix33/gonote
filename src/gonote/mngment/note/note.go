package note

import (
	"gonote/db"
	"html"
	"log"
	"time"

	"github.com/google/uuid"
)

// Note contains the content of a note.
type Note struct {
	ID      string
	Title   string
	Owner   string
	Added   time.Time
	Deleted bool
}

// Delete remove the note from the list.
// dbID is a reference to the database ID, if nil, a new connection will be opened.
func (n *Note) Delete(c *db.Conn) (e error) {
	db.MustConnect(c, func(c *db.Conn) {
		// Run delete query.
		query := `
			UPDATE Note
			SET Note.Deleted = 1
			WHERE Note.ID = ?`
		params := []interface{}{n.ID}
		_, _, e = db.Run(c, query, params, nil)

		if e != nil {
			log.Fatalln(e)
			return
		}
		log.Printf("DELETE %s(%s)\n", n.Title, n.ID)
	})

	return e
}

// Add takes a parital note and add it.
// If the note already exists, it clones it.
func (n *Note) Add(dbID *string) error {
	// Generate Object
	n.ID = uuid.New().String()
	n.Title = html.EscapeString(n.Title)
	n.Added = time.Now()

	var qErr error

	db.MustConnect(dbID, func(id string) {
		// Run Insert query.
		query := `
			INSERT INTO Note (
				ID, Title, Owner, Added
			) VALUES (?,?,?,?)`
		params := []interface{}{
			n.ID, n.Title, n.Owner, n.Added}
		_, _, qErr = db.Run(id, query, params, nil)

		if qErr != nil {
			log.Fatalln(qErr)
			return
		}
		log.Printf("ADD %s(%s)\n", n.Title, n.ID)
	})

	return qErr
}

// Update change the content and save to file.
func (n *Note) Update(n2 *Note, dbID *string) error {
	// Update the object.
	n.Title = html.EscapeString(n2.Title)

	var qErr error
	db.MustConnect(dbID, func(id string) {
		// Run update Query
		query := `
			UPDATE Note SET
				Note.Title = ?,
				Note.Owner = ?,
				Note.Added = ?,
				Note.Deleted = ?
			WHERE Note.ID = ?`
		params := []interface{}{
			n.Title, n.Owner, n.Added, n.Deleted, n.ID}
		_, _, qErr := db.Run(id, query, params, nil)

		if qErr != nil {
			log.Fatalln(qErr)
			return
		}
		log.Printf("UPDATE %s(%s)\n", n.Title, n.ID)
	})

	return qErr
}
