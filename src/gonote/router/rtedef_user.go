package router

import (
	"gonote/db"
	"gonote/mngment"
	"net/http"
)

const userRteAddr = "^/user/?{id}?$"

type userRteHandler struct{}

func (userRteHandler) Get(rw *http.ResponseWriter, req *http.Request, r *Route) {
	var id = r.Params["id"]

	if id != "" {
		// We can get a user from the database with the ID.

		// Get the User.
		var user *mngment.User
		db.MustConnect(nil, func(c *db.Conn) {
			user = mngment.GetUserByID(id, c)
		})

		if user != nil {
			// Secure the user details.
			if !r.User.IsAdmin && r.User.ID != user.ID {
				user.Password = ""
				user.Email = ""
			}

			WriteJSON(rw, user)
			return
		}
	} else if id == "" && r.User != nil {
		// We don't have the ID of a user, but we have an authenticated user.
		// Return the information about the current user.

		WriteJSON(rw, r.User)
		return
	}

	NotFound(rw)
}
