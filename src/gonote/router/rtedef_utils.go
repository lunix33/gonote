package router

import (
	"gonote/db"
	"gonote/mngment"
	"net/http"
)

const utilRteInfoAddr = "^/util/info$"

type utilRteInfoResponse struct {
	CustomStyle *string
	Title       *string
}

func utilRteInfo(rw *http.ResponseWriter, req *http.Request, r *Route) {
	// get the application information from the DB.
	var (
		customset *mngment.Setting
		titleset  *mngment.Setting
	)
	db.MustConnect(nil, func(c *db.Conn) {
		customset = mngment.GetSetting(mngment.CustomPathSetting, c)
		titleset = mngment.GetSetting(mngment.SiteTitleSetting, c)
	})

	// Make the response object.
	rsp := utilRteInfoResponse{}
	if customset != nil {
		rsp.CustomStyle = &(*customset).Value
	}
	if titleset != nil {
		rsp.Title = &(*titleset).Value
	}

	WriteJSON(rw, rsp)
}
