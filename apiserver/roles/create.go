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

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// CreateRoleRequest is the object that is expected by the
// Create() function.
type CreateRoleRequest struct {
	Name        string
	Description string `json:",omitempty"`
	Parameters  string `json:",omitempty"`
}

// Create creates a new Account object in
// the database.
func Create(r render.Render, params martini.Params, db database.Datastore, data CreateRoleRequest) {

	// Store the Role object in the database. In case the
	// database operation fails, an error response is sent back to the caller.
	newRole, err := DBCreateRole(db.Get(), data)
	if err != nil {
		responses.CreateError(r, err.Error())
		return
	}
	// Return the account.
	responses.CreateOK(r, newRole)
}

// DBCreateRole creates a new Account object in the database.
//
func DBCreateRole(db *sql.DB, data CreateRoleRequest) (RoleList, error) {

	stmt, err := db.Prepare(`INSERT platform_role SET name=?, description=?,
		 parameters=?`)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(data.Name, data.Description, data.Parameters)
	if err != nil {
		return nil, err
	}

	// The id of the newly generated Role
	RoleID, _ := res.LastInsertId()
	// Retrieve the newly created object from the database and return it
	Roles, err := DBGetRoles(db, RoleID)
	if err != nil {
		return nil, err
	}
	// Return the account object.
	return Roles, nil
}
