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

// Delete removes one or more Role objects from the database and
// sends them back to caller.
//
func Delete(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// RoleID is either -1 if no Role ID was provided or > 0 otherwise.
	var RoleID int64 = -1

	// Convert the Role ID string to a 64-bit integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		RoleID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.DeleteError(r, fmt.Sprintf("Invalid Role ID: %v", RoleID))
			return
		}
	}

	// Delete the Role object from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	err := DBDeleteRole(db.Get(), RoleID)
	if err != nil {
		responses.DeleteError(r, err.Error())
		return
	}

	// Return the modified Role.
	responses.DeleteOK(r, "Role deleted")
}

// DBDeleteRole removes the Role from the MySQL database.
func DBDeleteRole(db *sql.DB, RoleID int64) error {

	stmt, err := db.Prepare(`
		DELETE from platform_role WHERE id=?`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(RoleID)
	if err != nil {
		return err
	}

	return nil
}
