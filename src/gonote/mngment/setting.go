package mngment

import (
	"gonote/db"
	"reflect"
)

// GetSetting retreive the keys
// `key` is the setting key to be fetched.
// `c` is an optional database connection
// Returns the setting (s) found.
func GetSetting(key string, c *db.Conn) (s *Setting) {
	s = new(Setting)
	s.Key = key

	db.MustConnect(c, func(c *db.Conn) {
		params := []interface{}{key}
		rst, count, err := db.Run(c, settingGetQuery, params, reflect.TypeOf(Setting{}))
		if err != nil {
			s = nil
		} else if count > 0 {
			rsetting := rst[0].(*Setting)
			s = rsetting
		}
	})

	return s
}

// GetAllSettings get all the application settings.
// `c` is an optional database connection.
// Returns a map with the setting key as a key and the setting object as a value for all the settings.
func GetAllSettings(c *db.Conn) (s map[string]*Setting) {
	s = make(map[string]*Setting)

	db.MustConnect(c, func(c *db.Conn) {
		rst, _, err := db.Run(c, settingGetAllQuery, nil, reflect.TypeOf(Setting{}))
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
	// PortSetting gives the setting key for the application port.
	PortSetting = "Port"

	// DBVersionSetting gives the setting key for the database version.
	DBVersionSetting = "DBVersion"

	// CustomPathSetting gives the setting key for the custom css path.
	CustomPathSetting = "CustomPath"

	// InterfaceSetting gives the setting key for the application interface.
	InterfaceSetting = "Interface"
)

// Setting represent a key-value pair of setting.
type Setting struct {
	Key   string
	Value string
}

// Set change the value of the setting.
// `v` is the value of the setting.
// `c` is an optional database connection
// Returns any error (e) occured.
func (s *Setting) Set(v string, c *db.Conn) (e error) {
	s.Value = v

	db.MustConnect(c, func(c *db.Conn) {
		var (
			q string
			p []interface{}
		)

		if v != "" {
			// If the value isn't nil, insert the value (or update if already present).
			q = settingUpsertQuery
			p = []interface{}{s.Key, s.Value, s.Value}
		} else {
			// Or delete the setting key.
			q = settingDeleteQuery
			p = []interface{}{s.Key}
		}
		_, _, err := db.Run(c, q, p, nil)
		if err != nil {
			e = err
		}
	})

	return e
}
