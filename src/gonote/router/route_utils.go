package router

import (
	"encoding/base64"
	"errors"
	"gonote/db"
	"gonote/mngment"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// isFile allow to find out if the path of a request is a file.
//
// "req" is the request object.
//
// Returns true if the request path is a file, otherwise false.
func isFile(req *http.Request) bool {
	match, matchErr := regexp.MatchString("[^.]+\\.[^.]+$", req.URL.Path)
	if matchErr != nil {
		return false
	}
	return match
}

// getParams finds the parameters of a path.
//
// "matcher" is the path pattern used to find the parameters.
// "req" is the request object.
//
// Returns:
// (p) The a map with the parameters.
// (e) Any error occured.
func getParams(matcher string, req *http.Request) (p map[string]string, e error) {
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

// getBody read the body of a request
//
// "req" is the request object.
//
// Returns
// (b) The byte slice representing the request body.
// (e) Any error occured.
func getBody(req *http.Request) (b []byte, e error) {
	defer func() {
		if r := recover(); r != nil {
			b = make([]byte, 0)
			e = errors.New("unable to get request body")
		}
	}()

	body, err := ioutil.ReadAll(req.Body)
	return body, err
}

// getContentType find the mimetype of a path.
//
// "path" is the path from which the mimetype should be detected.
//
// Returns a string with the mimetype representation.
func getContentType(path string) string {
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

// decodeToken decodes the authentication token sent in a request header.
//
// "req" is the request object.
//
// Returns
// (u) The username.
// (t) The user token.
// (e) Any error occured.
func decodeToken(req *http.Request) (u string, t string, e error) {
	if req.Header["Authorization"] == nil {
		return u, t, errors.New("no authorization header")
	}

	headerValue := req.Header["Authorization"][0]
	signature := headerValue[6:len(headerValue)]
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return u, t, err
	}

	split := strings.Split(string(decoded), ":")
	if len(split) > 1 {
		u = split[0]
		t = split[1]
		return u, t, nil
	}

	return u, t, errors.New("unable to get signature information")
}

// Authenticate gets the user sending the request.
//
// "req" is the request object.
// "c" is an optional database connection.
//
// Returns the user (u) logged in.
func Authenticate(req *http.Request, c *db.Conn) (u *mngment.User) {
	username, token, err := decodeToken(req)
	if err != nil {
		return
	}

	db.MustConnect(c, func(c *db.Conn) {
		ut := mngment.GetUserToken(token, username, c)
		if ut != nil {
			u = mngment.GetUserByID(ut.UserID, c)
		}
	})

	return
}
