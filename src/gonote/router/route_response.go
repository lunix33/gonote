package router

import (
	"encoding/json"
	"fmt"
	"gonote/mngment"
	"log"
	"net/http"
)

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
	Error   HTTPStatusResponseError
}

// HTTPStatusResponseError is the error message emitted to the client.
type HTTPStatusResponseError struct {
	Message  string `json:"message"`
	Stack    string `json:"stack"`
	Friendly string `json:"friendly"`
}

// NotFound respond to the client request with a 404 error (not found).
//
// "rw" is the object used to respond to the client request.
func NotFound(rw *http.ResponseWriter) {
	status := HTTPStatus{
		Status: HTTPStatusNotFound}
	WriteResponse(rw, status)
}

// Unauthorized respond to the client request with a 400 error (unauthorized).
//
// "rw is the object used to respond to the client request."
func Unauthorized(rw *http.ResponseWriter) {
	status := HTTPStatus{
		Status: HTTPUnauthorized}
	WriteResponse(rw, status)
}

// InternalError respond to the client request with a 500 error (internal error).
//
// "rw" is the object used to respond to the client request.
// "err" is an error object to send with the response.
// "friendly" is an alternative error message made for regular users.
// "usr" is the user logged in.
func InternalError(rw *http.ResponseWriter, err error, friendly string, usr *mngment.User) {
	// Make json error.
	status := HTTPStatus{
		Status: HTTPStatusError}

	// Add the error object based on the user logged in.
	if usr.IsAdmin {
		status.Error = HTTPStatusResponseError{
			Message: fmt.Sprintf("%v", err),
			Stack:   fmt.Sprintf("%+v", err)}
	} else {
		status.Error = HTTPStatusResponseError{
			Friendly: friendly}
	}

	WriteResponse(rw, status)
}

// WriteJSON Send to the client a JSON representation of an object.
//
// "rw" is the object used to respond to the client request.
// "obj" is the object to be included as a json message.
func WriteJSON(rw *http.ResponseWriter, obj interface{}) {
	status := HTTPStatus{
		Status:  HTTPStatusOk,
		Message: obj}
	WriteResponse(rw, status)
}

// WriteResponse Write the response message to the client.
//
// "rw" is the object used to respond to the client request.
// "status" is the response status object.
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
