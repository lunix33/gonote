package db

import (
	"database/sql"
	"errors"
	"fmt"
	"gonote/util"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/google/uuid"
)

var (
	connections = make(map[string]*sql.DB)
	dbFile      = util.DirnameJoin("notes.db")
	box         = packr.New("sql", util.DirnameJoin("db", "sql"))
)

// Connect create a connection with the database file.
// Returns	(i) The id of the database.
//			(e) Any occurred error.
func Connect() (i string, e error) {
	dbc, oerr := sql.Open("sqlite3", dbFile)
	if oerr != nil {
		return i, oerr
	}

	i = uuid.New().String()
	connections[i] = dbc

	if _, e = os.Stat(dbFile); os.IsNotExist(e) {
		initDb(i)
	}

	return i, nil
}

// Close terminate a connection with the database.
// `id` is the ID of the database connection.
// Returns an error if any occure.
func Close(id string) error {
	var (
		err error
		dbc = connections[id]
	)

	if id == "" {
		return nil
	} else if dbc != nil {
		err = dbc.Close()
	} else {
		err = errors.New("no connection with specified id")
	}

	return err
}

// MustConnect ensure a connection is made with the database before running some actions.
// `dbID` is a reference to the database ID, if nil a new connection will be created and terminated once the callback is done.
// `cb` is the callback function.
func MustConnect(dbID *string, cb func(id string)) {
	var id string
	if dbID == nil {
		cID, cErr := Connect()
		defer Close(cID)
		if cErr != nil {
			panic(cErr)
		}

		id = cID
	} else {
		id = *dbID
	}

	cb(id)
}

// Run execute a prepared query on a connected database.
// `id` is the database ID.
// `query` is the query string to be executed.
// `params` is the list of parameters in the query.
// `outType` is the type of the objects.
// Returns	(r) An array of type outType with the query results.
//			(c) Then the number of rows returned/affected.
//			(e) Any error occured.
func Run(id, query string, params []interface{}, outType reflect.Type) (r []interface{}, c int64, e error) {
	dbc := connections[id]
	if dbc != nil {
		// If the params list is nil, create an empty one.
		if params == nil {
			params = make([]interface{}, 0)
		}

		// Check which command should be ran.
		query = strings.TrimSpace(query)
		if strings.HasPrefix(query, "SELECT") && outType != nil {
			return queryQuery(dbc, query, params, outType)
		}
		return execQuery(dbc, query, params)
	}

	return r, c, errors.New("no connection with specified id")
}

// queryQuery run the query with Query (for SELECT)
// `dbc` is the database connection object.
// `query` is the query string to be executed.
// `params` is the list of parameters in the query.
// `outType` is the type of the objects.
// Returns	(r) An array of `outType` with the results of the query.
//			(c) The number of rows returned.
//			(e) Any error occured.
func queryQuery(dbc *sql.DB, query string, params []interface{}, outType reflect.Type) (r []interface{}, c int64, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = errors.New("unable to build request result")
		}
	}()

	var cols []string
	r = make([]interface{}, 0)

	// Run the query
	rows, queryErr := dbc.Query(query, params...)
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
				return nil, c, errCol
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
			return nil, c, scanErr
		}

		// Build the output struct.
		o := reflect.New(outType)
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
		r = append(r, o.Interface())

		c++
	}

	return r, c, nil
}

// execQuery run the query with Exec (for INSERT, UPDATE, DELETE, ...)
// `dbc` is the database connection object
// `query` is the query string to be executed.
// `params` is the list of parameters in the query.
// Returns	(r) An empty array...
//			(a) The number of affected rows.
//			(e) Any error occured.
func execQuery(dbc *sql.DB, query string, params []interface{}) (r []interface{}, a int64, e error) {
	r = make([]interface{}, 0)

	// Run the query
	res, queryErr := dbc.Exec(query, params...)
	if queryErr != nil {
		return nil, 0, queryErr
	}

	// Get the number of affected rows.
	affectedRow, affErr := res.RowsAffected()
	if affErr != nil {
		return nil, 0, affErr
	}

	return r, affectedRow, nil
}

// initDb initialize the database.
// `id` is the database ID.
func initDb(id string) {
	q, err := box.FindString("init.sql")
	if err != nil {
		panic(err)
	}
	_, _, qerr := Run(id, q, nil, nil)

	if qerr != nil {
		panic(qerr)
	} else {
		log.Println("Database properly initiated.")
	}
}

// MigrateFrom migrate the application database from one version to the latest.
// `version` is the version from which you are starting to migrate.
// `to` is the target version, if the value of `to` is 0, then run all the migrations.
// `dbID` is the ID of the database.
func MigrateFrom(version int64, to int64, dbID *string) {
	MustConnect(dbID, func(id string) {
		// Find the last migration number if `to` is 0.
		if to == 0 {
			to = findLastMigration(version)
		}

		// Apply each migration one at the time.
		for version++; version <= to; version++ {
			filename := fmt.Sprintf("migrations/m-%d.sql", version)

			// Get the content of the migration file
			q, err := box.FindString(filename)
			if err != nil {
				panic(err)
			}

			// Update query
			q = "BEGIN TRANSACTION;\n" + q
			q += `
				UPDATE Setting
				SET Value = ?
				WHERE Setting.Key = "DBVersion";
				COMMIT;`
			p := []interface{}{strconv.FormatInt(version, 10)}

			// Run migration file.
			_, _, qerr := Run(id, q, p, nil)
			if qerr != nil {
				log.Fatalln(qerr)
				panic(fmt.Sprintf("migration to db version %d failed.", version))
			} else {
				log.Printf("Migration to db version %d successful.", version)
			}
		}
	})
}

// findLastMigration tries to find the last available migration.
// `from` is the version from which we should start looking.
// Returns the last version number found.
func findLastMigration(from int64) int64 {
	from++
	filename := fmt.Sprintf("migrations/m-%d.sql", from)

	for box.Has(filename) {
		from++
		filename = fmt.Sprintf("migrations/m-%d.sql", from)
	}

	return from - 1
}
