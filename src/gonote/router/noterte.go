package router

import (
	"encoding/json"
	"gonote/db"
	"gonote/mngment"
	"net/http"
)

const noteRteSearchAddr = "/note/search"

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
		InternalError(rw, err)
		return
	}

	notes := make([]*mngment.Note, 0)
	db.MustConnect(nil, func(c *db.Conn) {
		notes = mngment.SearchNotes(crits, c)
	})

	WriteJSON(rw, notes)
}

const noteRteAddr = "/note/{id}"

//
func noteRte(rw *http.ResponseWriter, req *http.Request, r *Route) {
	if req.Method == http.MethodGet {
		noteRteGet(rw, req, r)
	} else if req.Method == http.MethodConnect {
		noteRteDelete(rw, req, r)
	} else if req.Method == http.MethodPut {
		noteRtePut(rw, req, r)
	}
}

//
func noteRteGet(rw *http.ResponseWriter, req *http.Request, r *Route) {
}

//
func noteRtePut(rw *http.ResponseWriter, req *http.Request, r *Route) {
}

//
func noteRteDelete(rw *http.ResponseWriter, req *http.Request, r *Route) {
}
