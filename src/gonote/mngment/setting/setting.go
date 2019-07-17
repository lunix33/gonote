package setting

import (
	"gonote/db"
	"reflect"
)

// Get retreive the keys
// `key` is the setting key to be fetched.
func Get(key string, dbID *string) (s *Setting) {
	s = new(Setting)
	s.Key = key

	db.MustConnect(dbID, func(id string) {
		q := `
			SELECT * FROM Setting
			WHERE Setting.Key = ?`
		params := []interface{}{key}
		rst, count, err := db.Run(id, q, params, reflect.TypeOf(Setting{}))
		if err != nil {
			s = nil
		} else if count > 0 {
			rsetting := rst[0].(*Setting)
			s = rsetting
		}
	})

	return s
}

// GetAll get all the application settings.
// `dbID` is the ID of the database.
// Returns a map with the setting key as a key and the setting object as a value for all the settings.
func GetAll(dbID *string) (s map[string]*Setting) {
	s = make(map[string]*Setting)

	db.MustConnect(dbID, func(id string) {
		q := `SELECT * FROM Setting`
		rst, _, err := db.Run(id, q, nil, reflect.TypeOf(Setting{}))
		if err == nil {
			for _, v := range rst {
				set := v.(*Setting)
				s[set.Key] = set
			}
		}
	})

	return s
}

const (
	// Port gives the setting key for the application port.
	Port = "Port"

	// DBVersion gives the setting key for the database version.
	DBVersion = "DBVersion"

	// CustomPath gives the setting key for the custom css path.
	CustomPath = "CustomPath"

	// Interface gives the setting key for the application interface.
	Interface = "Interface"
)

// Setting represent a key-value pair of setting.
type Setting struct {
	Key   string
	Value string
}

// Set change the value of the setting.
func (s *Setting) Set(v string, dbID *string) (e error) {
	s.Value = v

	db.MustConnect(dbID, func(id string) {
		var (
			q string
			p []interface{}
		)

		if v != "" {
			// If the value isn't nil, insert the value (or update if already present).
			q = `
				INSERT INTO Setting(Key, Value) VALUES (?, ?)
				ON CONFLICT(Key) DO UPDATE
				SET Value = ?`
			p = []interface{}{s.Key, s.Value, s.Value}
		} else {
			// Or delete the setting key.
			q = `
				DELETE FROM Setting
				WHERE Setting.Key = ?`
			p = []interface{}{s.Key}
		}
		_, _, err := db.Run(id, q, p, nil)
		if err != nil {
			e = err
		}
	})

	return e
}
