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

// globalHandler is the general request handler.
//
// "next" is the specialized function called once this function is done.
//
// Returns a handler function usable by the http lib.
func globalHandler(rw http.ResponseWriter, req *http.Request) {
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

// serveDefault is the default route handler.
//
// "rw" is the response object.
// "req" is the request object.
// "r" is the route detail.
func serveDefault(rw *http.ResponseWriter, req *http.Request, r *Route) {
	var (
		isf      = isFile(req)
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
	mimetype = getContentType(path)

	(*rw).Header().Set("Content-Type", mimetype)
	(*rw).WriteHeader(http.StatusOK)
	(*rw).Write(content)
}

// RegisterRoute register the HTTP routes of the application.
func RegisterRoute() {
	routes[securityRteLoginAddr] = securityRteLogin
	routes[securityRteLogoutAddr] = securityRteLogout
	routes[noteRteSearchAddr] = noteRteSearch
	routes[noteRteAddr] = noteRte

	http.HandleFunc("/", globalHandler)
}
