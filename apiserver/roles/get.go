//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package roles

import (
	"net/http"

	"github.com/oweidner/platform/apiserver/responses"
	"github.com/oweidner/platform/apiserver/utils"
	"github.com/oweidner/platform/database"

	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
)

// List returns the list of available plans.
//
func List(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {
	// Retreive the (list of) plans from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	var roles RoleList

	_, err := db.GetDBMap().Select(&roles, "SELECT * FROM platform_role WHERE _deleted=0 ORDER BY id")
	if err != nil {
		responses.Error(r, err.Error())
		return
	}

	responses.OKStatusPlusData(r, roles, len(roles))
}

// Get retrieves one or more plan objects from the database and
// sends them back to caller.
//
func Get(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// Parse the resource ID into an int64
	resourceID, parseError := utils.ParseResourceID(params["p1"])
	if parseError != nil {
		responses.Error(r, parseError.Error())
	}

	// Retreive the (list of) plans from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	var role Role
	err := db.GetDBMap().SelectOne(&role, "SELECT * FROM platform_role WHERE _deleted=0 AND id=?", resourceID)
	if err != nil {
		responses.Error(r, err.Error())
		return
	}

	responses.OKStatusPlusData(r, role, 1)
}
