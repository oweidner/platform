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

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// CreateOrgRequest is the object that is expected by the
// Create() function.
type CreateOrgRequest struct {
	Orgname string
	Name    string
	Email   string
}

// Create creates a new organiztion object in
// the database.
func Create(r render.Render, params martini.Params, db database.Datastore, data CreateOrgRequest) {

	// Store the organiztion object in the database. In case the
	// database operation fails, an error response is sent back to the caller.
	newOrg, err := DBCreateOrg(db.Get(), data)
	if err != nil {
		responses.CreateError(r, err.Error())
		return
	}
	// Return the organiztion.
	responses.CreateOK(r, newOrg)
}

// DBCreateOrg creates a new Account object in the database.
//
func DBCreateOrg(db *sql.DB, data CreateOrgRequest) (OrgList, error) {

	stmt, err := db.Prepare(`INSERT platform_organisation SET orgname=?, name=?, email=?`)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(data.Orgname, data.Name, data.Email)
	if err != nil {
		return nil, err
	}

	// The id of the newly generated organiztion
	orgID, _ := res.LastInsertId()
	// Retrieve the newly created object from the database and return it
	orgs, err := DBGetOrgs(db, orgID)
	if err != nil {
		return nil, err
	}
	// Return the organiztion object.
	return orgs, nil
}
