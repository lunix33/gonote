package mngment

import (
	"fmt"
	"gonote/db"
	"reflect"
	"strings"
)

// NoteSearchCriterions are the search criterions for notes.
type NoteSearchCriterions struct {
	Username *string
	Trash    *string
	Public   *string
	Text     *string
	Order    *string
	Limit    *int
	Offset   *int
}

// GetNote retrive a note with a specified id.
// If the note ain't found, returns nil.
func GetNote(nID string, c *db.Conn) (n *Note) {
	db.MustConnect(c, func(c *db.Conn) {
		p := []interface{}{nID}
		rst, cnt, err := db.Run(c, noteGetQuery, p, reflect.TypeOf(Note{}))
		if err == nil && cnt > 0 {
			n = rst[0].(*Note)
		}
	})

	return n
}

// SearchNotes fetch all the notes from the DB which correspond to criterias.
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
		}
	})

	return sr
}

func noteBuildWhere(crits NoteSearchCriterions) (w string, p []interface{}) {
	w = "WHERE "
	p = make([]interface{}, 0, 3)
	clauses := make([]string, 0, 4)

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

	w += strings.Join(clauses, " AND\n")

	return w, p
}
