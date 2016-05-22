//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package self

import (
	"net/http"

	"github.com/oweidner/platform/apiserver/accounts"
	"github.com/oweidner/platform/apiserver/authentication"
	"github.com/oweidner/platform/apiserver/responses"
	"github.com/oweidner/platform/database"
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
)

// GetSelf retrieves the account referenced in the auth token.
//
func GetSelf(req *http.Request, params martini.Params, r render.Render, db database.Datastore, user authentication.UserInfo) {

	// account holds the data returned to the caller
	var account accounts.Account

	// Query the database for the given UserID
	err := db.GetDBMap().SelectOne(&account, "SELECT * FROM platform_account WHERE _deleted=0 AND id=?", user.UserID)
	if err != nil {
		responses.Error(r, err.Error())
		return
	}

	// "gray-out" the password
	account.Password = ""

	responses.OKStatusPlusData(r, account, 1)
}
