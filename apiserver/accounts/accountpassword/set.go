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

	"github.com/oweidner/platform/apiserver/authentication"
	"github.com/oweidner/platform/apiserver/responses"
	"github.com/oweidner/platform/database"
	"golang.org/x/crypto/bcrypt"

	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
)

// Create inserts a new object in the database.
func Set(r render.Render, params martini.Params, db database.Datastore, data AccountPassword, user authentication.UserInfo) {

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
	_, err := db.GetDB().Exec("UPDATE platform_account SET password=? WHERE id=?", pwdHash, user.UserID)
	if err != nil {
		responses.Error(r, err.Error())
		return
	}

	responses.OKStatusOnly(r, "Password Changed")

}
