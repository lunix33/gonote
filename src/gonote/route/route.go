package route

import (
	"fmt"
	"gonote/note"
	"io/ioutil"
	"net/http"
)

const (
	defaultPath            = "/"
	apiNotesPath           = "/api/notes/"
	apiNotesPathWithParams = "/api/notes/{author}"
)

// *: *
// (200) Return the index file.
// *: *.*
// (200) Return the requested file in the public folder.
func serveDefault(rw http.ResponseWriter, req *http.Request) {
	var (
		isf  = IsFile(req)
		path string
	)

	// If the request is a file get the path to the requested file in public,
	// Otherwise use 'public/index.html'
	if path = "public/index.html"; isf {
		path = fmt.Sprintf("public%s", req.URL.Path)
	}

	// Read file.
	fc, err := ioutil.ReadFile(path)
	if err != nil {
		InternalError(rw, req, err)
		return
	}

	// Get content type
	ct := GetContentType(path)

	rw.Header().Set("Content-Type", ct)
	rw.WriteHeader(http.StatusOK)
	rw.Write(fc)
}

// GET: /api/notes
// (404) Not found.
// (200:XHR) Return the list of notes.
func apiNotes(rw http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		params, err := GetParams(apiNotesPathWithParams, req)
		if err != nil {
			InternalError(rw, req, err)
			return
		}

		lst := note.List(params["author"])
		WriteJSON(rw, req, lst)
	} else {
		NotFound(rw, req)
	}
}

// RegisterRoute register the HTTP routes of the application.
func RegisterRoute() {
	http.Handle(defaultPath, Cors(serveDefault))
	http.Handle(apiNotesPath, Cors(apiNotes))
	http.Handle(apiNotePath, Cors(apiNote))
}
