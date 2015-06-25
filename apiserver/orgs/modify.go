//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package orgs

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// ModifyOrgRequest is the object that is expected by the
// Modify() function.
type ModifyOrgRequest struct {
	Name string `binding:"required"`
}

// Modify modifies a org object in the database.
//
func Modify(r render.Render, params martini.Params, db database.Datastore, data ModifyOrgRequest) {

	// orgID is either -1 if no orgm ID was provided or > 0 otherwise.
	var orgID int64 = -1

	// Convert the org ID string to a 64-bit integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		orgID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.ModifyError(r, fmt.Sprintf("Invalid Organization ID: %v", orgID))
			return
		}
	}
	// Update the org object in the database. In case the
	// database operation fails, an error response is sent back to the caller.
	modifiedOrg, err := DBModifyOrg(db.Get(), orgID, data)
	if err != nil {
		responses.ModifyError(r, err.Error())
		return
	}

	// Return the modified org.
	responses.ModifyOK(r, modifiedOrg)
}

// DBModifyOrg modifies a Account object in the database.
//
func DBModifyOrg(db *sql.DB, orgID int64, data ModifyOrgRequest) (OrgList, error) {

	stmt, err := db.Prepare(`
		UPDATE organization SET name=? WHERE id=?`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(data.Name, orgID)
	if err != nil {
		return nil, err
	}

	// Retrieve the modified object from the database and return it
	account, err := DBGetOrgs(db, orgID)
	if err != nil {
		return nil, err
	}

	return account, nil
}
