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
	"fmt"
	"strconv"

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// ModifyAccountRequest is the object that is expected by the
// Modify() function.
type ModifyAccountRequest struct {
	Firstname string `binding:"required"`
	Lastname  string `binding:"required"`
	Email     string `binding:"required"`
	Password  string `binding:"required"`
}

// Modify modifies a account object in the database.
//
func Modify(r render.Render, params martini.Params, db database.Datastore, data ModifyAccountRequest) {

	// accountID is either -1 if no account ID was provided or > 0 otherwise.
	var accountID int64 = -1

	// Convert the accountid string to an integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		accountID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.ModifyError(r, fmt.Sprintf("Invalid Account ID: %v", accountID))
			return
		}
	}
	// Retrieve the (list of) account from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	modifiedAccount, err := DBModifyAccount(db.Get(), accountID, data)
	if err != nil {
		responses.GetError(r, err.Error())
		return
	}

	// Return the account.
	responses.ModifyOK(r, modifiedAccount)
}

// DBModifyAccount modifies a Account object in the database.
//
func DBModifyAccount(db *sql.DB, accountID int64, data ModifyAccountRequest) (AccountList, error) {

	stmt, err := db.Prepare(`
		UPDATE account SET firstname=?, lastname=?, email=?, password=?
        WHERE id=?`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(
		data.Firstname, data.Lastname, data.Email, data.Password, accountID)
	if err != nil {
		return nil, err
	}

	// Retrieve the modified object from the database and return it
	account, err := DBGetAccounts(db, accountID)
	if err != nil {
		return nil, err
	}

	return account, nil
}
