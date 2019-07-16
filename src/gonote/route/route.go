package route

import (
	"fmt"
	"gonote/models/note"
	"gonote/util"
	"net/http"
	"os"

	"github.com/gobuffalo/packr/v2"
)

const (
	defaultPath            = "/"
	apiNotesPath           = "/api/notes/"
	apiNotesPathWithParams = "/api/notes/{author}"
)

var box = packr.New("builtin", util.DirnameJoin("builtin"))

// GlobalHandler is the general request handler.
// `next` is the specialized function called once this function is done.
// Returns a handler function usable by the http lib.
func GlobalHandler(next func(*http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			// Global error handling.
			defer func() {
				if r := recover(); r != nil {
					InternalError(&rw, fmt.Errorf("%v", r))
				}
			}()

			// Set CORS headers.
			rw.Header().Set("Access-Control-Allow-Methods", "GET, PUT, PATCH, DELETE, OPTIONS")
			rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			rw.Header().Set("Access-Control-Allow-Origin", "*")

			// Calling next function if not OPTIONS request.
			if req.Method == http.MethodOptions {
				rw.Write([]byte(""))
			} else {
				next(&rw, req)
			}
		})
}

// *: *
// (200) Return the index file.
// *: *.*
// (200) Return the requested file in the public folder.
func serveDefault(rw *http.ResponseWriter, req *http.Request) {
	var (
		isf      = IsFile(req)
		path     string
		content  []byte
		mimetype string
		err      error
	)

	// If the request is a file get the path to the requested file in public,
	// Otherwise use 'builtin/index.html'
	if path = fmt.Sprintf("builtin%s", req.URL.Path); !isf {
		path = "index.html"
	}

	// Get the file from the builtin box
	content, err = box.Find(path)
	if err != nil && os.IsNotExist(err) {
		NotFound(rw)
		return
	} else if err != nil {
		InternalError(rw, err)
		return
	}
	mimetype = GetContentType(path)

	(*rw).Header().Set("Content-Type", mimetype)
	(*rw).WriteHeader(http.StatusOK)
	(*rw).Write(content)
}

// GET: /api/notes
// (404) Not found.
// (200:XHR) Return the list of notes.
func apiNotes(rw *http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		params, err := GetParams(apiNotesPathWithParams, req)
		if err != nil {
			InternalError(rw, err)
			return
		}

		lst := note.List(params["author"], nil)
		WriteJSON(rw, lst)
	} else {
		NotFound(rw)
	}
}

// RegisterRoute register the HTTP routes of the application.
func RegisterRoute() {
	http.Handle(defaultPath, GlobalHandler(serveDefault))
	http.Handle(apiNotesPath, GlobalHandler(apiNotes))
	http.Handle(apiNotePath, GlobalHandler(apiNote))
}
