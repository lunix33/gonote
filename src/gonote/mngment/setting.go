package mngment

import (
	"gonote/db"
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
