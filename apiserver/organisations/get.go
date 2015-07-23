//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package organisations

import (
	"net/http"

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/apiserver/utils"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// List returns the list of available resources.
//
func List(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	db.GetDBMap().AddTableWithName(Organisation{}, "platform_organisation").SetKeys(true, "id")

	// Retreive the requested resource from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	var organisations OrganisationList
	_, err := db.GetDBMap().Select(&organisations, "SELECT * FROM platform_organisation ORDER BY id")
	if err != nil {
		responses.Error(r, err.Error())
		return
	}
	responses.OKStatusPlusData(r, organisations, len(organisations))
}

// Get retrieves one or more resource objects from the database and
// sends them back to caller.
//
func Get(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	db.GetDBMap().AddTableWithName(Organisation{}, "platform_account").SetKeys(true, "id")

	// Parse the resource ID into an int64
	resourceID, parseError := utils.ParseResourceID(params["p1"])
	if parseError != nil {
		responses.Error(r, parseError.Error())
	}

	// Retreive the list of all resources from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	var organisation Organisation
	err := db.GetDBMap().SelectOne(&organisation, "SELECT * FROM platform_organisation where id=?", resourceID)
	if err != nil {
		responses.Error(r, err.Error())
		return
	}
	responses.OKStatusPlusData(r, organisation, 1)
}