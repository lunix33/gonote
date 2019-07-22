package router

import (
	"fmt"
	"gonote/util"
	"log"
	"net/http"
	"os"

	"github.com/gobuffalo/packr/v2"
)

var (
	box    = packr.New("builtin", util.DirnameJoin("builtin"))
	routes = make(map[string]RouteFn)
)

// GlobalHandler is the general request handler.
// `next` is the specialized function called once this function is done.
// Returns a handler function usable by the http lib.
func GlobalHandler(rw http.ResponseWriter, req *http.Request) {
	// Global error handling.
	defer func() {
		if r := recover(); r != nil {
			InternalError(&rw, fmt.Errorf("%v", r))
		}
	}()

	// Set CORS headers.
	rw.Header().Set("Access-Control-Allow-Methods", "GET, PUT, PATCH, DELETE, OPTIONS")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// Calling next function if not OPTIONS request.
	if req.Method == http.MethodOptions {
		rw.Write([]byte(""))
	} else {
		route := findRoute(req)
		if route == nil {
			route = &Route{
				Handler: serveDefault,
				Matcher: "/"}
		}

		log.Printf("%s: %s", req.Method, route.Matcher)
		route.Handler(&rw, req, route)
	}
}

// *: *
// (200) Return the index file.
// *: *.*
// (200) Return the requested file in the public folder.
func serveDefault(rw *http.ResponseWriter, req *http.Request, r *Route) {
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

func findRoute(req *http.Request) (r *Route) {
	var (
		err    error
		params map[string]string
	)

	// Find route handler and params.
	for k, v := range routes {
		params, err = GetParams(k, req)
		if err == nil {
			r = &Route{
				Params:  params,
				Matcher: k,
				Handler: v}

			if req.Method == http.MethodPatch ||
				req.Method == http.MethodPost ||
				req.Method == http.MethodPut {
				r.Body, err = GetBody(req)
			}
			break
		}
	}

	return r
}

// RegisterRoute register the HTTP routes of the application.
func RegisterRoute() {

	http.HandleFunc("/", GlobalHandler)
}
