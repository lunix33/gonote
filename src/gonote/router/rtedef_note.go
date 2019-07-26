package router

import (
	"encoding/json"
	"gonote/db"
	"gonote/mngment"
	"net/http"
	"time"
)

const noteRteSearchAddr = "^/note/search$"

type noteRteSearchResponse struct {
	Note       *mngment.Note
	LastUpdate time.Time
	Author     string
}

// noteRteGet respond to the "/note/search" (POST) route.
// It gets all the notes which respect the search filters.
func noteRteSearch(rw *http.ResponseWriter, req *http.Request, r *Route) {
	if req.Method != http.MethodPost {
		NotFound(rw)
		return
	}

	crits := mngment.NoteSearchCriterions{}
	err := json.Unmarshal(r.Body, &crits)
	if err != nil {
		InternalError(rw, err, "We were unable to decifer your search criterions.", r.User)
		return
	}

	notes := make([]*mngment.Note, 0)
	rtn := make([]*noteRteSearchResponse, 0)
	db.MustConnect(nil, func(c *db.Conn) {
		// Forcing some search criterions for security reason if ...
		// The user isn't an admin
		// And isn't searching for its own notes.
		if r.User == nil || (!r.User.IsAdmin && *crits.Username != r.User.Username) {
			*crits.Public = "only"
			crits.Trash = nil
		}

		notes = mngment.SearchNotes(crits, c)

		for _, v := range notes {
			user := mngment.GetUserByID(v.UserID, c)
			contents := v.GetNoteContent(c)
			ro := &noteRteSearchResponse{
				Note:       v,
				LastUpdate: (*contents[0]).Updated,
				Author:     user.Username}
			rtn = append(rtn, ro)
		}
	})

	WriteJSON(rw, rtn)
}

const noteRteAddr = "^/note/{id}"

// noteRteHandler is a virtual struct to store the routes functions.
type noteRteHandler struct{}

// Get respond to the "/note/{id}" (GET) routes.
// It retreives a note from the database.
func (noteRteHandler) Get(rw *http.ResponseWriter, req *http.Request, r *Route) {

}

// Delete respond to the "/note/{id}" (DELETE) routes.
// It trashs or removes a note from the database.
func (noteRteHandler) Delete(rw *http.ResponseWriter, req *http.Request, r *Route) {

}

// Put respond to the "/note/{id}" (PUT) routes.
// It updates a note in the database.
func (noteRteHandler) Put(rw *http.ResponseWriter, req *http.Request, r *Route) {

}
