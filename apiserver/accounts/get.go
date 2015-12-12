//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package accounts

import (
	"net/http"

	"github.com/oweidner/platform/apiserver/responses"
	"github.com/oweidner/platform/apiserver/utils"
	"github.com/oweidner/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// List returns the list of available resources.
//
func List(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// Retreive the requested resource from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	var accounts AccountList
	_, err := db.GetDBMap().Select(&accounts, "SELECT * FROM platform_account WHERE _deleted=0 ORDER BY id")
	if err != nil {
		responses.Error(r, err.Error())
		return
	}

	// "black-out" the passwords
	for idx := range accounts {
		accounts[idx].Password = "*****"
	}

	responses.OKStatusPlusData(r, accounts, len(accounts))
}

// Get retrieves one or more resource objects from the database and
// sends them back to caller.
//
func Get(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// Parse the resource ID into an int64
	resourceID, parseError := utils.ParseResourceID(params["p1"])
	if parseError != nil {
		responses.Error(r, parseError.Error())
	}

	// Retreive the list of all resources from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	var account Account
	err := db.GetDBMap().SelectOne(&account, "SELECT * FROM platform_account WHERE _deleted=0 AND id=?", resourceID)
	if err != nil {
		responses.Error(r, err.Error())
		return
	}

	// "black-out" the password
	account.Password = "***"

	responses.OKStatusPlusData(r, account, 1)
}
