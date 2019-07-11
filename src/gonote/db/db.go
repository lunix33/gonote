package db

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

const dbFile = "notes.db"

var connections = make(map[string]*sql.DB)

// Connect create a connection with the database file.
// The function retuns the connection ID once properly created.
func Connect() (string, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return "", err
	}

	id := uuid.New().String()
	connections[id] = db

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		initDb(id)
	}

	return id, nil
}

// Close terminate a connection with the database.
// It takes in the ID of the connection to close.
func Close(id string) error {
	var (
		err error
		db  = connections[id]
	)

	if id == "" {
		return nil
	} else if db != nil {
		err = db.Close()
	} else {
		err = errors.New("no connection with specified id")
	}

	return err
}

// Run execute a prepared query on a connected database.
func Run(id, query string, params []interface{}, outType reflect.Type) ([]interface{}, int64, error) {
	db := connections[id]
	if db != nil {
		// If the params list is nil, create an empty one.
		if params == nil {
			params = make([]interface{}, 0)
		}

		// Check which command should be ran.
		query = strings.TrimSpace(query)
		if strings.HasPrefix(query, "SELECT") && outType != nil {
			return queryQuery(db, query, params, outType)
		}
		return execQuery(db, query, params)
	}

	return nil, 0, errors.New("no connection with specified id")
}

// queryQuery run the query with Query (for SELECT)
func queryQuery(db *sql.DB, query string, params []interface{}, outType reflect.Type) ([]interface{}, int64, error) {
	var (
		affected int64
		rtn      = make([]interface{}, 0)
		cols     []string
	)

	// Run the query
	rows, queryErr := db.Query(query, params...)
	if queryErr != nil {
		return nil, 0, queryErr
	}

	// Fetch the results.
	defer rows.Close()
	for rows.Next() {
		// Get the list of cols to reflect in the struct.
		if cols == nil {
			colsRow, errCol := rows.Columns()
			if errCol != nil {
				return nil, affected, errCol
			}
			cols = colsRow
		}

		// Fetch the column into an array.
		binder := make([]interface{}, len(cols))
		for i := range binder {
			binder[i] = new(interface{})
		}
		scanErr := rows.Scan(binder...)
		if scanErr != nil {
			return nil, affected, scanErr
		}

		// Build the output struct.
		o := reflect.New(outType).Elem()
		for i, v := range cols {
			field := reflect.Indirect(o).FieldByName(v)
			value := *(binder[i].(*interface{}))
			valueType := reflect.TypeOf(value)

			if valueType.Kind() == reflect.Int {
				field.SetInt(value.(int64))
			} else if valueType.Kind() == reflect.String {
				field.SetString(value.(string))
			} else if valueType.Kind() == reflect.Bool {
				field.SetBool(value.(bool))
			} else if valueType.Kind() == reflect.Float64 {
				field.SetFloat(value.(float64))
			}
		}
		rtn = append(rtn, o.Interface())

		affected++
	}

	return rtn, affected, nil
}

// execQuery run the query with Exec (for INSERT, UPDATE, DELETE, ...)
func execQuery(db *sql.DB, query string, params []interface{}) ([]interface{}, int64, error) {
	// Run the query
	res, queryErr := db.Exec(query, params...)
	if queryErr != nil {
		return nil, 0, queryErr
	}

	// Get the number of affected rows.
	affectedRow, affErr := res.RowsAffected()
	if affErr != nil {
		return nil, 0, affErr
	}

	return make([]interface{}, 0), affectedRow, nil
}

func initDb(id string) {
	q := `
		CREATE TABLE notes (
			ID TEXT PRIMARY KEY,
			Title TEXT,
			Author TEXT,
			Content TEXT,
			Added TEXT,
			Updated TEXT
		)`
	_, _, qerr := Run(id, q, nil, nil)

	if qerr != nil {
		log.Fatalln(qerr)
	} else {
		log.Println("Database properly initiated.")
	}
}
