//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package accountpassword

import (
	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/apiserver/utils"
	"github.com/codewerft/platform/database"
	"golang.org/x/crypto/bcrypt"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// Create inserts a new object in the database.
func Set(r render.Render, params martini.Params, db database.Datastore, data AccountPassword) {

	// Parse the resource ID into an int64
	resourceID, parseError := utils.ParseResourceID(params["p1"])
	if parseError != nil {
		responses.Error(r, parseError.Error())
	}

	if len(data.Password.String) < 8 {
		responses.Error(r, "Passwords need to be at least 8 characters long.")
		return
	}

	// Create a bcrypt hash from the password as we don't want to store
	// plain-text passwords in the database
	pwdHash, bcryptError := bcrypt.GenerateFromPassword([]byte(data.Password.String), 0)
	if bcryptError != nil {
		responses.Error(r, bcryptError.Error())
		return
	}

	// Update the database record
	stmt, err := db.GetDB().Prepare("UPDATE platform_account SET password=? WHERE id=?")
	if err != nil {
		responses.Error(r, err.Error())
		return
	}
	_, err = stmt.Exec(pwdHash, resourceID)
	if err != nil {
		responses.Error(r, err.Error())
		return
	}

	responses.OKStatusOnly(r, "Password Changed")

}
