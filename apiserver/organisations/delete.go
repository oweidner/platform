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

// Delete removes one or more plan objects from the database and
// sends them back to caller.
//
func Delete(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	db.GetDBMap().AddTableWithName(Organisation{}, "platform_organisation").SetKeys(true, "id")

	// Parse the resource ID into an int64
	resourceID, parseError := utils.ParseResourceID(params["p1"])
	if parseError != nil {
		responses.Error(r, parseError.Error())
	}

	// Delete the object from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	_, err := db.GetDBMap().Exec("DELETE FROM platform_organisation WHERE id=?", resourceID)
	if err != nil {
		responses.Error(r, err.Error())
		return
	}

	// Return the modified plan.
	responses.OKStatusOnly(r, "Resource deleted")
}
