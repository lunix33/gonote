package router

import (
	"encoding/json"
	"gonote/db"
	"gonote/mngment"
	"log"
	"net/http"
	"strings"
)

type loginPostData struct {
	Username string
	Password string
}

const securityRteLoginAddr = "/login"

// securityRteLogin reponds to request to the "/login" (POST) route.
// It logins the user and send back the appropriate login token.
func securityRteLogin(rw *http.ResponseWriter, req *http.Request, r *Route) {
	if req.Method != http.MethodPost {
		NotFound(rw)
		return
	}

	postData := loginPostData{}
	err := json.Unmarshal(r.Body, &postData)
	if err != nil {
		InternalError(rw, err)
		return
	}

	db.MustConnect(nil, func(c *db.Conn) {
		// Validate the user identity.
		user := mngment.GetUser(postData.Username, nil)
		if user != nil && user.ComparePassword(postData.Password) {
			// Create new token for the user.
			ipSlice := strings.Split(req.RemoteAddr, ":")
			tok := mngment.UserToken{
				UserID: user.ID,
				IP:     ipSlice[0]}
			err = tok.Add(c)
			if err != nil {
				InternalError(rw, err)
				return
			}

			// Login successful
			log.Printf("-> Login %s (%s)", user.Username, tok.IP)
			WriteJSON(rw, tok)
			return
		}

		// The identity doesn't match.
		Unauthorized(rw)
	})
}

const securityRteLogoutAddr = "/logout"

// securityRteLogout respond to the "/logout" (GET) route.
// It logouts the user from the system by deleting the user token.
func securityRteLogout(rw *http.ResponseWriter, req *http.Request, r *Route) {
	if req.Method != http.MethodGet {
		NotFound(rw)
		return
	}

	username, token, err := decodeToken(req)
	if err == nil {
		db.MustConnect(nil, func(c *db.Conn) {
			// Get the token from the database.
			user := mngment.GetUser(username, c)
			tok := mngment.GetUserToken(token, user.ID, c)

			if tok != nil {
				// Delete the token.
				err = tok.Delete(c)
				if err != nil {
					InternalError(rw, err)
					return
				}

				// Token deleted.
				WriteJSON(rw, nil)
				return
			}

			// Token was not found.
			NotFound(rw)
		})
	}
}
