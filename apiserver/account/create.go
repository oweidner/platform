//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package account

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	"github.com/codewerft/platform/apiserver/responses"

	"github.com/codewerft/platform/database"
	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// CreateAccountRequest is the object that is expected by the
// CreateAccount() function.
type CreateAccountRequest struct {
	Firstname string `binding:"required"`
	Lastname  string `binding:"required"`
	Email     string `binding:"required"`
	Username  string `binding:"required"`
	Password  string `binding:"required"`
}

// Create creates a new Account object in
// the database.
func Create(r render.Render, params martini.Params, db database.Datastore, data CreateAccountRequest) {

	// Store the account object the database. In case the
	// database operation fails, an error response is sent back to the caller.
	newAccount, err := DBCreateAccount(db.Get(), data)
	if err != nil {
		responses.GetError(r, err.Error())
		return
	}
	// Return the account.
	responses.CreateOK(r, newAccount)
}

// DBCreateAccount creates a new Account object in the database.
//
func DBCreateAccount(db *sql.DB, data CreateAccountRequest) (AccountList, error) {

	// Create a bcrypt hash from the password as we don't want to store
	// plain-text passwords in the database
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 0)
	if err != nil {
		// TODO: handle properly
	}

	stmt, err := db.Prepare(`
		INSERT platform_account SET firstname=?, lastname=?, email=?,
		username=?,password=?`)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(
		data.Firstname, data.Lastname, data.Email, data.Username, pwdHash)
	if err != nil {
		return nil, err
	}

	// The id of the newly generated account
	accountID, _ := res.LastInsertId()
	// Retrieve the newly created object from the database and return it
	account, err := DBGetAccounts(db, accountID)
	if err != nil {
		return nil, err
	}
	// Return the account object.
	return account, nil
}
