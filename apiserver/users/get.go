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
	"net/http"
	"strconv"

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// Get retrieves one or more account objects from the database and
// sends them back to caller.
//
func Get(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// accountID is either -1 if no account ID was provided or > 0 otherwise.
	var accountID int64 = -1

	// Convert the accountid string to an integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		accountID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.GetError(r, fmt.Sprintf("Invalid Account ID: %v", accountID))
			return
		}
	}
	// Retrieve the (list of) account from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	account, err := DBGetAccounts(db.Get(), accountID)
	if err != nil {
		responses.GetError(r, err.Error())
		return
	}

	// Return the list of account or a 404 if the account wasn't found.
	if accountID != -1 && len(account) < 1 {
		responses.GetNotFound(r)
	} else {
		responses.GetOK(r, account)
	}
}

// DBGetAccounts returns a Account object from the database.
func DBGetAccounts(db *sql.DB, accountID int64) (AccountList, error) {

	// If no accountID is provided (accountID is -1), all account are retreived. If
	// a accountID is given, a specific account is retreived.
	var rows *sql.Rows

	if accountID == -1 {
		queryString := `SELECT * FROM platform_account`
		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query()
		if err != nil {
			return nil, err
		}

	} else {
		queryString := `SELECT * FROM platform_account WHERE id = ?`
		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query(accountID)
		if err != nil {
			return nil, err
		}
	}

	// Read the rows into the target struct
	var objs AccountList

	for rows.Next() {

		var obj Account
		err := rows.Scan(
			&obj.ID, &obj.Firstname, &obj.Lastname, &obj.Email,
			&obj.Username, &obj.Password)

		// Forward the error
		if err != nil {
			return nil, err
		}
		// Append object to the list
		objs = append(objs, obj)
	}

	return objs, nil
}
