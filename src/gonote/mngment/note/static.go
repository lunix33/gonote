package note

import (
	"gonote/db"
	"log"
	"reflect"
)

// Get retrive a note with a specified id.
// If the note ain't found, returns nil.
func Get(nID string, dbID *string) *Note {
	// Connect to DB
	id, cErr := db.Connect()
	defer db.Close(id)
	if cErr != nil {
		log.Fatalln(cErr)
		return nil
	}

	var rst *Note
	db.MustConnect(dbID, func(id string) {
		query := `
			SELECT * FROM notes WHERE ID = ?`
		params := []interface{}{nID}
		qRst, count, qErr := db.Run(id, query, params, reflect.TypeOf(Note{}))
		if qErr != nil {
			log.Fatalln(qErr)
			return
		} else if count > 0 {
			// Return the first result.
			n := qRst[0].(Note)
			rst = &n
		}
	})

	return rst
}

// List fetch all the notes from the DB and send it back.
func List(author string, dbID *string) []*Note {
	list := make([]*Note, 0)
	db.MustConnect(dbID, func(id string) {
		query := `SELECT * FROM notes WHERE Author = ?`
		params := []interface{}{author}
		rst, _, qErr := db.Run(id, query, params, reflect.TypeOf(Note{}))
		if qErr != nil {
			log.Fatalln(qErr)
			return
		}

		for i := range rst {
			n := rst[i].(Note)
			list = append(list, &n)
		}
	})

	return list
}
