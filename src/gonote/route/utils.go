package route

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

// IsFile allow to find out if the path of a request is a file.
func IsFile(req *http.Request) bool {
	match, matchErr := regexp.MatchString("[^.]+\\.[^.]+$", req.URL.Path)
	if matchErr != nil {
		return false
	}
	return match
}

// GetParams finds the parameters of a path.
// The "matcher" is the path pattern used to find the parameters.
func GetParams(matcher string, req *http.Request) (map[string]string, error) {
	rtn := make(map[string]string)

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
				rtn[v] = matches[0][i+1]
			}
		}
	}

	return rtn, err
}

// GetBody read the body of a request and returns it.
func GetBody(req *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	return body, err
}

// GetContentType find the content type of a path.
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

// Return code for a message.
const (
	HTTPStatusOk       = "ok"
	HTTPStatusNotFound = "not_found"
	HTTPStatusError    = "error"
)

// HTTPStatus is the global response message.
type HTTPStatus struct {
	Status  string
	Message interface{}
	Error   error
}

// Cors setup requests to handle CORS request.
func Cors(next func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("Access-Control-Allow-Methods", "GET, PUT, PATCH, DELETE, OPTIONS")
			rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			if req.Method == http.MethodOptions {
				rw.Write([]byte(""))
			} else {
				next(rw, req)
			}
		})
}

// NotFound return a not found (404) error to the client.
func NotFound(rw http.ResponseWriter, req *http.Request) {
	status := HTTPStatus{
		Status: HTTPStatusNotFound}
	WriteResponse(rw, req, status)
}

// InternalError send through the responseWriter an internal error.
func InternalError(rw http.ResponseWriter, req *http.Request, err error) {
	// Make json error.
	status := HTTPStatus{
		Status: HTTPStatusError,
		Error:  err}
	WriteResponse(rw, req, status)
}

// WriteJSON Send to the client a JSON representation of an object.
func WriteJSON(rw http.ResponseWriter, req *http.Request, obj interface{}) {
	status := HTTPStatus{
		Status:  HTTPStatusOk,
		Message: obj}
	WriteResponse(rw, req, status)
}

// WriteResponse Write the response message to the client.
func WriteResponse(rw http.ResponseWriter, req *http.Request, status HTTPStatus) {
	jsonObj, jsonErr := json.Marshal(status)
	if jsonErr != nil {
		log.Fatalln("Unable to write response.")
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	switch status.Status {
	case HTTPStatusNotFound:
		rw.WriteHeader(http.StatusNotFound)
	case HTTPStatusError:
		rw.WriteHeader(http.StatusInternalServerError)
	default:
		rw.WriteHeader(http.StatusOK)
	}
	rw.Write(jsonObj)
}
