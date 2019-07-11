package note

import (
	"fmt"
	"gonote/db"
	"log"
	"reflect"
	"time"
)

// Get retrive a note with a specified id.
// If the note ain't found, returns nil.
func Get(nID string) *Note {
	// Connect to DB
	id, cErr := db.Connect()
	defer db.Close(id)
	if cErr != nil {
		log.Fatalln(cErr)
		return nil
	}

	// Run Select query
	query := `
		SELECT * FROM notes WHERE ID = ?`
	params := []interface{}{nID}
	rst, count, qErr := db.Run(id, query, params, reflect.TypeOf(Note{}))
	if qErr != nil {
		log.Fatalln(qErr)
		return nil
	} else if count > 0 {
		// Return the first result.
		n := rst[0].(Note)
		return &n
	}
	// If no result.
	return nil
}

// List fetch all the notes from the DB and send it back.
func List(author string) []*Note {
	// Connect to DB
	id, cErr := db.Connect()
	defer db.Close(id)
	if cErr != nil {
		log.Fatalln(cErr)
		return nil
	}

	// Run Select Query
	query := `SELECT * FROM notes WHERE Author = ?`
	params := []interface{}{author}
	rst, _, qErr := db.Run(id, query, params, reflect.TypeOf(Note{}))
	if qErr != nil {
		log.Fatalln(qErr)
		return nil
	}

	list := make([]*Note, len(rst))
	for i := range rst {
		n := rst[i].(Note)
		list[i] = &n
	}
	return list
}

func dateString() string {
	tn := time.Now()
	date := fmt.Sprintf(
		"%d/%d/%d %d:%d:%d",
		tn.Year(), tn.Month(), tn.Day(),
		tn.Hour(), tn.Minute(), tn.Second())

	return date
}
