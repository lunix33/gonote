package route

import (
	"encoding/json"
	"gonote/models/note"
	"net/http"
)

const (
	apiNotePath       = "/api/note/"
	apiNoteWithIDPath = "/api/note/{id}"
)

func apiNote(rw *http.ResponseWriter, req *http.Request) {
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
		NotFound(rw)
	}
}

// GET: /api/note/:ID:
func getAPINote(rw *http.ResponseWriter, req *http.Request) {
	// Get the url params.
	params, paramsErr := GetParams(apiNoteWithIDPath, req)
	if paramsErr != nil {
		InternalError(rw, paramsErr)
		return
	}

	// Get the item.
	noteItem := note.Get(params["id"], nil)
	if noteItem == nil {
		NotFound(rw)
		return
	}

	WriteJSON(rw, noteItem)
}

// PUT: /api/note
func putAPINote(rw *http.ResponseWriter, req *http.Request) {
	// Get the request body.
	body, err := GetBody(req)
	if err != nil {
		InternalError(rw, err)
		return
	}

	// Decode the body into a new note.
	noteData := new(note.Note)
	json.Unmarshal(body, noteData)

	// Add the new note.
	addErr := noteData.Add(nil)
	if addErr != nil {
		InternalError(rw, addErr)
		return
	}

	WriteJSON(rw, noteData)
}

// PATCH: /api/note/:ID:
func patchAPINote(rw *http.ResponseWriter, req *http.Request) {
	// Get params
	params, paramsErr := GetParams(apiNoteWithIDPath, req)
	if paramsErr != nil {
		InternalError(rw, paramsErr)
		return
	}

	// Get body
	body, bodyErr := GetBody(req)
	if bodyErr != nil {
		InternalError(rw, bodyErr)
		return
	}

	// Parse body to new note.
	noteUp := new(note.Note)
	jsonErr := json.Unmarshal(body, noteUp)
	if jsonErr != nil {
		InternalError(rw, jsonErr)
		return
	}

	// Get current note.
	noteEle := note.Get(params["id"], nil)
	if noteEle == nil {
		NotFound(rw)
		return
	}

	// Update current note with new note.
	updateErr := noteEle.Update(noteUp, nil)
	if updateErr != nil {
		InternalError(rw, updateErr)
		return
	}

	WriteJSON(rw, noteEle)
}

// DELETE: /api/note/:ID:
func deleteAPINote(rw *http.ResponseWriter, req *http.Request) {
	// Get the params
	params, paramsErr := GetParams(apiNoteWithIDPath, req)
	if paramsErr != nil {
		InternalError(rw, paramsErr)
		return
	}

	// Get the note.
	noteItem := note.Get(params["id"], nil)
	if noteItem == nil {
		NotFound(rw)
		return
	}

	// Delete the item
	err := noteItem.Delete(nil)
	if err != nil {
		InternalError(rw, err)
		return
	}

	WriteJSON(rw, nil)
}
