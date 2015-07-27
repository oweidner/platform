//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package accountplans

import (
	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/apiserver/utils"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// Create inserts a new object in the database.
func Create(r render.Render, params martini.Params, db database.Datastore, data AccountPlanAssoc) {

	db.GetDBMap().AddTableWithName(AccountPlanAssoc{}, "platform_account_plan_assoc").SetKeys(true, "id")

	// Parse the resource ID into an int64
	accountID, parseError := utils.ParseResourceID(params["p1"])
	if parseError != nil {
		responses.Error(r, parseError.Error())
	}

	if accountID != data.AccountID.Int64 {
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
