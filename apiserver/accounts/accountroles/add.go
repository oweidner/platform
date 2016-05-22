//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package accountroles

import (
	"github.com/oweidner/platform/apiserver/responses"
	"github.com/oweidner/platform/apiserver/utils"
	"github.com/oweidner/platform/database"

	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
)

// Create inserts a new object in the database.
func Create(r render.Render, params martini.Params, db database.Datastore, data AccountOrganisationRole) {

	// Parse the resource ID into an int64
	accountID, parseError := utils.ParseResourceID(params["p1"])
	if parseError != nil {
		responses.Error(r, parseError.Error())
	}

	if accountID != data.AccountID {
		responses.Error(r, "accountID != data.AccountID")
		return
	}

	// Store the object in the database. In case the
	// database operation fails, an error response is sent back to the caller.
	err := db.GetDBMap().Insert(&data)
	if err != nil {
		responses.Error(r, err.Error())
		return
	}
	responses.OKStatusPlusData(r, data, 1)
}
