package router

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

// IsFile allow to find out if the path of a request is a file.
//
// "req" is the request object.
//
// Returns true if the request path is a file, otherwise false.
func IsFile(req *http.Request) bool {
	match, matchErr := regexp.MatchString("[^.]+\\.[^.]+$", req.URL.Path)
	if matchErr != nil {
		return false
	}
	return match
}

// GetParams finds the parameters of a path.
//
// "matcher" is the path pattern used to find the parameters.
// "req" is the request object.
//
// Returns:
// (p) The a map with the parameters.
// (e) Any error occured.
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

		pathRegexStr := reg.ReplaceAllString(matcher, `([^/]*)`)
		reg, err = regexp.Compile(pathRegexStr)
		if err == nil {
			matches = reg.FindAllStringSubmatch(req.URL.Path, -1)
			if len(matches) > 0 {
				// Makes sure the url was matched.
				p = make(map[string]string)
				for i, v := range paramNames {
					p[v] = matches[0][i+1]
				}
			} else {
				err = errors.New("matcher unable to match url")
			}
		}
	}

	return p, err
}

// GetBody read the body of a request
//
// `req` is the request object.
//
// Returns
// (b) The byte slice representing the request body.
// (e) Any error occured.
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
//
// `path` is the path from which the mimetype should be detected.
//
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
