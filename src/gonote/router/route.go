package router

import (
	"fmt"
	"gonote/mngment"
	"gonote/util"
	"log"
	"net/http"
	"os"

	"github.com/gobuffalo/packr/v2"
)

var (
	box    = packr.New("builtin", util.DirnameJoin("builtin"))
	routes = make([]RoutePair, 5)
)

// globalHandler is the general request handler.
//
// "next" is the specialized function called once this function is done.
//
// Returns a handler function usable by the http lib.
func globalHandler(rw http.ResponseWriter, req *http.Request) {
	var route *Route

	// Global error handling.
	defer func() {
		if r := recover().(error); r != nil {
			var u *mngment.User
			if route != nil {
				u = route.User
			}
			InternalError(&rw, r, "An error occured and makes us unable to continue.", u)
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
		route = findRoute(req)
		if route == nil {
			route = &Route{
				Handler: serveDefault,
				Matcher: "/"}
		} else if route.Handler == nil {
			NotFound(&rw)
			return
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
		InternalError(rw, err, "An error occured while trying to get the requested resource.", r.User)
		return
	}
	mimetype = getContentType(path)

	(*rw).Header().Set("Content-Type", mimetype)
	(*rw).WriteHeader(http.StatusOK)
	(*rw).Write(content)
}

// RegisterRoute register the HTTP routes of the application.
func RegisterRoute() {
	// Initialize route handlers (routes with multiple methods.)
	noteRte := noteRteHandler{}

	// Route registration
	routes[0] = RoutePair{
		Path:     securityRteLoginAddr,
		Handlers: MethodHandler{http.MethodPost: securityRteLogin}}

	routes[1] = RoutePair{
		Path:     securityRteLogoutAddr,
		Handlers: MethodHandler{http.MethodGet: securityRteLogout}}

	routes[2] = RoutePair{
		Path:     noteRteSearchAddr,
		Handlers: MethodHandler{http.MethodPost: noteRteSearch}}

	routes[3] = RoutePair{
		Path: noteRteAddr,
		Handlers: MethodHandler{
			http.MethodGet:    noteRte.Get,
			http.MethodDelete: noteRte.Delete,
			http.MethodPut:    noteRte.Put}}

	routes[4] = RoutePair{
		Path:     utilRteInfoAddr,
		Handlers: MethodHandler{http.MethodGet: utilRteInfo}}

	// Register the router with the http server.
	http.HandleFunc("/", globalHandler)
}
