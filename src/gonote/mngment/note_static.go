package mngment

import (
	"fmt"
	"gonote/db"
	"gonote/util"
	"reflect"
	"strings"
)

// NoteSearchCriterions are the search criterions for notes.
type NoteSearchCriterions struct {
	Username *string
	Trash    *string
	Public   *string
	Text     *string
	Tags     *[]string
	DateFrom *string
	DateTo   *string
	Order    *string
	Limit    *int
	Offset   *int
}

// GetNote retrive a note with a specified id.
//
// "nID" is the ID of the note.
// "c" is an optional database connection.
//
// Returns the note (n) with the specified ID, or nil.
func GetNote(nID string, c *db.Conn) (n *Note) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{nID}
		rst, cnt, err := db.Run(c, noteGetQuery, p, reflect.TypeOf(Note{}))
		if err == nil && cnt > 0 {
			n = rst[0].(*Note)
		} else if err != nil {
			util.LogErr(err)
		}
	})

	return n
}

// SearchNotes fetch all the notes from the DB which correspond to criterias.
//
// "crits" are the search criterions.
// "c" is an optional database connection
//
// Returns all the notes (sr) which correspond to the search criterions.
func SearchNotes(crits NoteSearchCriterions, c *db.Conn) (sr []*Note) {
	db.MustConnect(c, func(c *db.Conn) {
		// Build WHERE
		where, p := noteBuildWhere(crits)

		// Order
		if crits.Order != nil {
			var field string
			switch *crits.Order {
			case "updated":
				field = `"NoteContent"."Updated"`
			case "added":
				field = `"Note"."Added"`
			case "version":
				field = `"NoteContent"."Version"`
			case "user":
				field = `"User"."Username"`
			default:
				field = ""
			}

			if field != "" {
				where += "\nORDER BY " + field
			}
		}

		// Limit
		if crits.Limit != nil {
			p = append(p, crits.Limit)
			where += "\nLIMIT ?"
		}

		// Offset
		if crits.Offset != nil {
			p = append(p, crits.Offset)
			where += "\nOFFSET ?"
		}

		// Run query
		q := fmt.Sprintf(`%s %s`, noteSearchQueryBase, where)
		rst, _, err := db.Run(c, q, p, reflect.TypeOf(Note{}))
		if err == nil {
			for _, v := range rst {
				sr = append(sr, v.(*Note))
			}
		} else {
			util.LogErr(err)
		}
	})

	return sr
}

// noteBuildWhere allow to build there where statement for the SearchNotes query.
//
// "crits" are the search criterions
//
// Returns:
// (w) The where statement.
// (p) The parameters associated with the statement.
func noteBuildWhere(crits NoteSearchCriterions) (w string, p []interface{}) {
	w = "WHERE "
	p = make([]interface{}, 0, 5)
	clauses := make([]string, 0, 6)

	// User search
	if crits.Username != nil {
		// If a user id is supplied.
		username := fmt.Sprintf("%%%s%%", *crits.Username)
		p = append(p, username)
		clauses = append(clauses, `"User"."Username" LIKE ?`)
	}

	// Trash search
	if crits.Trash != nil && *crits.Trash == "only" {
		clauses = append(clauses, `"Note"."Deleted" = 1`)
	} else if crits.Trash != nil && *crits.Trash == "include" {
		clauses = append(clauses, `"Note"."Deleted" <= 1`)
	} else {
		clauses = append(clauses, `"Note"."Deleted" = 0`)
	}

	// Public search
	if crits.Public != nil && *crits.Public == "only" {
		clauses = append(clauses, `"Note"."Public" = 1`)
	} else if crits.Public != nil && *crits.Public == "exclude" {
		clauses = append(clauses, `"Note"."Public" = 0`)
	}

	// Text search
	if crits.Text != nil {
		text := fmt.Sprintf("%%%s%%", *crits.Text)
		p = append(p, text, text)
		clauses = append(clauses, `(
			"Note"."Title" LIKE ? OR
			"NoteContent"."Content" LIKE ?
		)`)
	}

	// Tag search
	if crits.Tags != nil {
		// Add all the tags to the parameter array.
		for _, v := range *crits.Tags {
			p = append(p, v)
		}

		paramsStr := ""
		if l := len(*crits.Tags); l > 0 {
			paramsStr = strings.Repeat("?,", l)[:l-1]
		}
		clauses = append(clauses, fmt.Sprintf(
			`"NoteTag"."Name" IN (%s)`, paramsStr))
	}

	// Date search
	if crits.DateFrom != nil && crits.DateTo != nil {
		p = append(p, crits.DateFrom, crits.DateTo)
		clauses = append(clauses, `"NoteContent"."Updated" BETWEEN ? AND ?`)
	}

	w += strings.Join(clauses, " AND\n")

	return w, p
}
