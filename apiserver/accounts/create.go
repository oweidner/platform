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
	"github.com/oweidner/platform/apiserver/responses"
	"github.com/oweidner/platform/database"
	"golang.org/x/crypto/bcrypt"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// Create inserts a new object in the database.
func Create(r render.Render, params martini.Params, db database.Datastore, data Account) {

	// Create a bcrypt hash from the password as we don't want to store
	// plain-text passwords in the database
	pwdHash, bcryptError := bcrypt.GenerateFromPassword([]byte(data.Password), 0)
	if bcryptError != nil {
		responses.Error(r, bcryptError.Error())
	}

	// Set the hashed password
	data.Password = string(pwdHash)

	// Store the object in the database. In case the
	// database operation fails, an error response is sent back to the caller.
	err := db.GetDBMap().Insert(&data)
	if err != nil {
		responses.Error(r, err.Error())
		return
	}
	responses.OKStatusPlusData(r, data, 1)
}
