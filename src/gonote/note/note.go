package note

import (
	"gonote/db"
	"html"
	"log"

	"github.com/google/uuid"
)

const noteFile = "notes.json"

// Note contains the content of a note.
type Note struct {
	ID      string
	Title   string
	Author  string
	Content string
	Added   string
	Updated string
}

// Delete remove the note from the list.
func (n *Note) Delete() error {
	// Connect to the DB.
	id, cErr := db.Connect()
	defer db.Close(id)
	if cErr != nil {
		log.Fatalln(cErr)
		return cErr
	}

	// Run delete query.
	query := `DELETE FROM notes WHERE ID=?`
	params := []interface{}{n.ID}
	_, _, qErr := db.Run(id, query, params, nil)
	
	if (qErr != nil) {
		log.Fatalln(qErr)
	}
	log.Printf("DELETE %s(%s)\n", n.Title, n.ID)

	return qErr
}

// Add takes a parital note and add it.
// If the note already exists, it clones it.
func (n *Note) Add() error {
	// Generate Object
	date := dateString()
	n.ID = uuid.New().String()
	n.Title = html.EscapeString(n.Title)
	n.Content = html.EscapeString(n.Content)
	n.Author = html.EscapeString(n.Author)
	n.Added = date
	n.Updated = date

	// Connect DB
	id, cErr := db.Connect()
	defer db.Close(id)
	if cErr != nil {
		log.Fatalln(cErr)
		return cErr
	}

	// Run Insert query.
	query := `
		INSERT INTO notes(
			ID, Title, Author, Content, Added, Updated
		) VALUES (?,?,?,?,?,?)`
	params := []interface{}{
		n.ID, n.Title, n.Author, n.Content, n.Added, n.Updated}
	_, _, qErr := db.Run(id, query, params, nil)

	if (qErr != nil) {
		log.Fatalln(qErr)
	}
	log.Printf("ADD %s(%s)\n", n.Title, n.ID)

	return qErr
}

// Update change the content and save to file.
func (n *Note) Update(n2 *Note) error {
	// Update the object.
	date := dateString()
	n.Content = html.EscapeString(n2.Content)
	n.Title = html.EscapeString(n2.Title)
	n.Author = html.EscapeString(n2.Author)
	n.Updated = date

	// Connect to DB
	id, cErr := db.Connect()
	defer db.Close(id)
	if cErr != nil {
		log.Fatalln(cErr)
		return cErr
	}

	// Run update Query
	query := `
		UPDATE notes SET
			Title = ?, Author = ?, Content = ?, Updated = ?
		WHERE ID = ?`
	params := []interface{}{
		n.Title, n.Author, n.Content, n.Updated, n.ID}
	_, _, qErr := db.Run(id, query, params, nil)

	if (qErr != nil) {
		log.Fatalln(qErr)
	}
	log.Printf("UPDATE %s(%s)\n", n.Title, n.ID)

	return qErr
}
