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
