package router

import (
	"gonote/mngment"
	"net/http"
)

// RouteFn is a type of function made to be handle by the custom router.
type RouteFn func(*http.ResponseWriter, *http.Request, *Route)

// RoutePair make a link between a route path and its handlers.
type RoutePair struct {
	Path     string
	Handlers MethodHandler
}

// MethodHandler is a key-value association between the HTTP method and the route function.
type MethodHandler map[string]RouteFn

// Route is the definition of a route once it has beed processed by the custom router.
type Route struct {
	Matcher string
	Params  map[string]string
	Body    []byte
	Handler RouteFn
	User    *mngment.User
}

// findRoute find the route to execute from a list of route.
//
// "req" is the request object.
//
// Returns the route matching the request, or nil if no routes were found.
func findRoute(req *http.Request) (r *Route) {
	var (
		err    error
		params map[string]string
	)

	for _, v := range routes {
		// For each of the routes, try to get the params
		// if we can then it's the right route.
		params, err = getParams(v.Path, req)
		if err == nil {
			r = &Route{
				Params:  params,
				Matcher: v.Path,
				Handler: v.Handlers[req.Method],
				User:    Authenticate(req, nil)}

			// Get the request body, if proper method.
			if req.Method == http.MethodPatch ||
				req.Method == http.MethodPost ||
				req.Method == http.MethodPut {
				r.Body, _ = getBody(req)
			}
			break
		}
	}

	return r
}
