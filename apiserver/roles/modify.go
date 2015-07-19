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
	"strconv"

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// ModifyRoleRequest is the object that is expected by the
// Modify() function.
type ModifyRoleRequest struct {
	Name        string
	Description string `json:",omitempty"`
	Parameters  string `json:",omitempty"`
}

// Modify modifies a Role object in the database.
//
func Modify(r render.Render, params martini.Params, db database.Datastore, data ModifyRoleRequest) {

	// RoleID is either -1 if no Rolem ID was provided or > 0 otherwise.
	var RoleID int64 = -1

	// Convert the Role ID string to a 64-bit integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		RoleID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.ModifyError(r, fmt.Sprintf("Invalid Role ID: %v", RoleID))
			return
		}
	}
	// Update the Role object in the database. In case the
	// database operation fails, an error response is sent back to the caller.
	modifiedRole, err := DBModifyRole(db.Get(), RoleID, data)
	if err != nil {
		responses.ModifyError(r, err.Error())
		return
	}

	// Return the modified Role.
	responses.ModifyOK(r, modifiedRole)
}

// DBModifyRole modifies a Account object in the database.
//
func DBModifyRole(db *sql.DB, RoleID int64, data ModifyRoleRequest) (RoleList, error) {

	stmt, err := db.Prepare(`
		UPDATE platform_Role SET name=?, description=?, parameters=? WHERE id=?`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(data.Name, data.Description, data.Parameters, RoleID)
	if err != nil {
		return nil, err
	}

	// Retrieve the modified object from the database and return it
	account, err := DBGetRoles(db, RoleID)
	if err != nil {
		return nil, err
	}

	return account, nil
}
