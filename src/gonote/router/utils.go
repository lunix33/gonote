package router

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

// IsFile allow to find out if the path of a request is a file.
// `req` is the request object.
// Returns true if the request path is a file, otherwise false.
func IsFile(req *http.Request) bool {
	match, matchErr := regexp.MatchString("[^.]+\\.[^.]+$", req.URL.Path)
	if matchErr != nil {
		return false
	}
	return match
}

// GetParams finds the parameters of a path.
// `matcher` is the path pattern used to find the parameters.
// `req` is the request object.
// Returns	(p) The a map with the parameters.
//			(e) Any error occured.
func GetParams(matcher string, req *http.Request) (p map[string]string, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = errors.New("unable to get the parameters")
			p = make(map[string]string)
		}
	}()

	reg, err := regexp.Compile(`\{([^{}]+)\}`)
	if err == nil {
		matches := reg.FindAllStringSubmatch(matcher, -1)
		var paramNames []string
		for _, v := range matches {
			paramNames = append(paramNames, v[1])
		}

		pathRegexStr := reg.ReplaceAllString(matcher, `([^/]+)`)
		reg, err = regexp.Compile(pathRegexStr)
		if err == nil {
			matches = reg.FindAllStringSubmatch(req.URL.Path, -1)
			for i, v := range paramNames {
				p[v] = matches[0][i+1]
			}
		}
	}

	return p, err
}

// GetBody read the body of a request
// `req` is the request object.
// Returns	(b) The byte slice representing the request body.
//			(e) Any error occured.
func GetBody(req *http.Request) (b []byte, e error) {
	defer func() {
		if r := recover(); r != nil {
			b = make([]byte, 0)
			e = errors.New("unable to get request body")
		}
	}()

	body, err := ioutil.ReadAll(req.Body)
	return body, err
}

// GetContentType find the mimetype of a path.
// `path` is the path from which the mimetype should be detected.
// Returns a string with the mimetype representation.
func GetContentType(path string) string {
	var reg = regexp.MustCompile(`\.(\w+)$`)
	match := reg.FindAllStringSubmatch(path, -1)
	if match != nil {
		ext := match[0][1]
		switch ext {
		case "html":
			return "text/html"
		case "js":
			return "application/javascript"
		case "css":
			return "text/css"
		case "ico":
			return "image/x-icon"
		}
	}

	return "text/plain"
}

// RouteFn is a type of function made to be handle by the custom router.
type RouteFn func(*http.ResponseWriter, *http.Request, *Route)

// Route is the definition of a route once it has beed processed by the custom router.
type Route struct {
	Matcher string
	Params  map[string]string
	Body    []byte
	Handler RouteFn
}

// Return code for a message.
const (
	HTTPStatusOk       = "ok"
	HTTPStatusNotFound = "not_found"
	HTTPStatusError    = "error"
	HTTPUnauthorized   = "unauth"
)

// HTTPStatus is the global response message.
type HTTPStatus struct {
	Status  string
	Message interface{}
	Error   error
}

// NotFound respond to the client request with a 404 error (not found).
// `rw` is the object used to respond to the client request.
func NotFound(rw *http.ResponseWriter) {
	status := HTTPStatus{
		Status: HTTPStatusNotFound}
	WriteResponse(rw, status)
}

// InternalError respond to the client request with a 500 error (internal error).
func InternalError(rw *http.ResponseWriter, err error) {
	// Make json error.
	status := HTTPStatus{
		Status: HTTPStatusError,
		Error:  err}
	WriteResponse(rw, status)
}

// WriteJSON Send to the client a JSON representation of an object.
func WriteJSON(rw *http.ResponseWriter, obj interface{}) {
	status := HTTPStatus{
		Status:  HTTPStatusOk,
		Message: obj}
	WriteResponse(rw, status)
}

// WriteResponse Write the response message to the client.
func WriteResponse(rw *http.ResponseWriter, status HTTPStatus) {
	jsonObj, jsonErr := json.Marshal(status)
	if jsonErr != nil {
		log.Fatalln("Unable to write response.")
		return
	}

	(*rw).Header().Set("Content-Type", "application/json")
	switch status.Status {
	case HTTPStatusNotFound:
		(*rw).WriteHeader(http.StatusNotFound)
	case HTTPStatusError:
		(*rw).WriteHeader(http.StatusInternalServerError)
	case HTTPUnauthorized:
		(*rw).WriteHeader(http.StatusUnauthorized)
	default:
		(*rw).WriteHeader(http.StatusOK)
	}
	(*rw).Write(jsonObj)
}
