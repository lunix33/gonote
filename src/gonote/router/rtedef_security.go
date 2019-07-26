package router

import (
	"encoding/json"
	"gonote/db"
	"gonote/mngment"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type loginPostData struct {
	Username string
	Password string
}

type loginPostReturn struct {
	User  *mngment.User
	Token *mngment.UserToken
}

const securityRteLoginAddr = "^/login$"

// securityRteLogin reponds to request to the "/login" (POST) route.
// It logins the user and send back the appropriate login token.
func securityRteLogin(rw *http.ResponseWriter, req *http.Request, r *Route) {
	if req.Method != http.MethodPost {
		NotFound(rw)
		return
	}

	postData := loginPostData{}
	if err := errors.Wrap(json.Unmarshal(r.Body, &postData), "unable to parse request body"); err != nil {
		InternalError(rw, err, "We weren't able to get your login credentials correctly.", r.User)
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
			if err := errors.Wrap(tok.Add(c), "unable to add user token"); err != nil {
				InternalError(rw, err, "We weren't able to create your user session.", r.User)
				return
			}

			// Login successful
			log.Printf("-> Login %s (%s)", user.Username, tok.IP)
			rtn := loginPostReturn{
				User:  user,
				Token: &tok}
			WriteJSON(rw, rtn)
			return
		}

		// The identity doesn't match.
		Unauthorized(rw)
	})
}

const securityRteLogoutAddr = "^/logout$"

// securityRteLogout respond to the "/logout" (GET) route.
// It logouts the user from the system by deleting the user token.
func securityRteLogout(rw *http.ResponseWriter, req *http.Request, r *Route) {
	if req.Method != http.MethodGet {
		NotFound(rw)
		return
	}

	uid, token, err := decodeToken(req)
	if err == nil {
		db.MustConnect(nil, func(c *db.Conn) {
			// Get the token from the database.
			tok := mngment.GetUserToken(token, uid, c)

			if tok != nil {
				// Delete the token.
				if err = errors.Wrapf(tok.Delete(c), "unable to delete user token %s", tok.Token); err != nil {
					InternalError(rw, err, "We weren't able to delete your session.", r.User)
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
	InternalError(rw, errors.New("unable to decode request token"), "We weren't able to delete your session.", r.User)
}
