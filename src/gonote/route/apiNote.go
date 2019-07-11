package route

import (
	"encoding/json"
	"gonote/note"
	"net/http"
)

const (
	apiNotePath       = "/api/note/"
	apiNoteWithIDPath = "/api/note/{id}"
)

func apiNote(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getAPINote(rw, req)
	case http.MethodPut:
		putAPINote(rw, req)
	case http.MethodPatch:
		patchAPINote(rw, req)
	case http.MethodDelete:
		deleteAPINote(rw, req)
	default:
		NotFound(rw, req)
	}
}

// GET: /api/note/:ID:
func getAPINote(rw http.ResponseWriter, req *http.Request) {
	// Get the url params.
	params, paramsErr := GetParams(apiNoteWithIDPath, req)
	if paramsErr != nil {
		InternalError(rw, req, paramsErr)
		return
	}

	// Get the item.
	noteItem := note.Get(params["id"])
	if noteItem == nil {
		NotFound(rw, req)
		return
	}

	WriteJSON(rw, req, noteItem)
}

// PUT: /api/note
func putAPINote(rw http.ResponseWriter, req *http.Request) {
	// Get the request body.
	body, err := GetBody(req)
	if err != nil {
		InternalError(rw, req, err)
		return
	}

	// Decode the body into a new note.
	noteData := new(note.Note)
	json.Unmarshal(body, noteData)

	// Add the new note.
	addErr := noteData.Add()
	if addErr != nil {
		InternalError(rw, req, addErr)
		return
	}

	WriteJSON(rw, req, noteData)
}

// PATCH: /api/note/:ID:
func patchAPINote(rw http.ResponseWriter, req *http.Request) {
	// Get params
	params, paramsErr := GetParams(apiNoteWithIDPath, req)
	if paramsErr != nil {
		InternalError(rw, req, paramsErr)
		return
	}

	// Get body
	body, bodyErr := GetBody(req)
	if bodyErr != nil {
		InternalError(rw, req, bodyErr)
		return
	}

	// Parse body to new note.
	noteUp := new(note.Note)
	jsonErr := json.Unmarshal(body, noteUp)
	if jsonErr != nil {
		InternalError(rw, req, jsonErr)
		return
	}

	// Get current note.
	noteEle := note.Get(params["id"])
	if noteEle == nil {
		NotFound(rw, req)
		return
	}

	// Update current note with new note.
	updateErr := noteEle.Update(noteUp)
	if updateErr != nil {
		InternalError(rw, req, updateErr)
		return
	}

	WriteJSON(rw, req, noteEle)
}

// DELETE: /api/note/:ID:
func deleteAPINote(rw http.ResponseWriter, req *http.Request) {
	// Get the params
	params, paramsErr := GetParams(apiNoteWithIDPath, req)
	if paramsErr != nil {
		InternalError(rw, req, paramsErr)
		return
	}

	// Get the note.
	noteItem := note.Get(params["id"])
	if noteItem == nil {
		NotFound(rw, req)
		return
	}

	// Delete the item
	err := noteItem.Delete()
	if err != nil {
		InternalError(rw, req, err)
		return
	}

	WriteJSON(rw, req, nil)
}
