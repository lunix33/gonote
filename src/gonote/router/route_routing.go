package router

import "net/http"

// RouteFn is a type of function made to be handle by the custom router.
type RouteFn func(*http.ResponseWriter, *http.Request, *Route)

// RouteList is a list of routes with their MethodHandler.
type RouteList map[string]MethodHandler

// MethodHandler is a key-value association between the HTTP method and the route function.
type MethodHandler map[string]RouteFn

// Route is the definition of a route once it has beed processed by the custom router.
type Route struct {
	Matcher string
	Params  map[string]string
	Body    []byte
	Handler RouteFn
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

	for k, v := range routes {
		// For each of the routes, try to get the params
		// if we can then it's the right route.
		params, err = getParams(k, req)
		if err == nil {
			r = &Route{
				Params:  params,
				Matcher: k,
				Handler: v[req.Method]}

			// Get the request body, if proper method.
			if req.Method == http.MethodPatch ||
				req.Method == http.MethodPost ||
				req.Method == http.MethodPut {
				r.Body, err = getBody(req)
			}
			break
		}
	}

	return r
}
