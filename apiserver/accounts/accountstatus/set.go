//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package accountstatus

import (
	"github.com/oweidner/platform/apiserver/responses"
	"github.com/oweidner/platform/apiserver/utils"
	"github.com/oweidner/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// Set sets the status of the resource.
func Set(r render.Render, params martini.Params, db database.Datastore, data AccountStatus) {

	// Parse the resource ID into an int64
	resourceID, parseError := utils.ParseResourceID(params["p1"])
	if parseError != nil {
		responses.Error(r, parseError.Error())
	}

	// Update the database record
	stmt, err := db.GetDB().Prepare("UPDATE platform_account SET status_disabled=? WHERE id=?")
	if err != nil {
		responses.Error(r, err.Error())
		return
	}
	_, err = stmt.Exec(data.Disabled, resourceID)
	if err != nil {
		responses.Error(r, err.Error())
		return
	}

	responses.OKStatusOnly(r, "Status Changed")

}
