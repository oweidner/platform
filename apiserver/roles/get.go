//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package roles

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

// Get retrieves one or more Role objects from the database and
// sends them back to caller.
//
func Get(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// RoleID is either -1 if no Role ID was provided or > 0 otherwise.
	var RoleID int64 = -1

	// Convert the Role ID string to an integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		RoleID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.GetError(r, fmt.Sprintf("Invalid Role ID: %v", RoleID))
			return
		}
	}
	// Retrieve the (list of) Roles from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	Roles, err := DBGetRoles(db.Get(), RoleID)
	if err != nil {
		responses.GetError(r, err.Error())
		return
	}

	// Return the list of Roles or a 404 if the account wasn't found.
	if RoleID != -1 && len(Roles) < 1 {
		responses.GetNotFound(r)
	} else {
		responses.GetOK(r, Roles)
	}
}

// DBGetRoles returns a Role object from the database.
func DBGetRoles(db *sql.DB, RoleID int64) (RoleList, error) {

	// If no accountID is provided (accountID is -1), all account are retreived. If
	// a RoleID is given, a specific account is retreived.
	var rows *sql.Rows

	if RoleID == -1 {
		queryString := `SELECT id, name, description, parameters
      FROM platform_role`

		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query()
		if err != nil {
			return nil, err
		}

	} else {
		queryString := `SELECT id, name, description, parameters
      FROM platform_Role WHERE id = ?`

		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query(RoleID)
		if err != nil {
			return nil, err
		}
	}

	// Read the rows into the target struct
	var objs RoleList

	for rows.Next() {

		var obj Role
		err := rows.Scan(
			&obj.ID, &obj.Name, &obj.Description, &obj.Parameters)

		// Forward the error
		if err != nil {
			return nil, err
		}
		// Append object to the list
		objs = append(objs, obj)
	}

	return objs, nil
}
