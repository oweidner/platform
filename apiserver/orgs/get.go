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
	"net/http"
	"strconv"

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// Get retrieves one or more organization objects from the database and
// sends them back to caller.
//
func Get(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// orgID is either -1 if no organization ID was provided or > 0 otherwise.
	var orgID int64 = -1

	// Convert the organization ID string to an integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		orgID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.GetError(r, fmt.Sprintf("Invalid Organization ID: %v", orgID))
			return
		}
	}
	// Retrieve the (list of) organizations from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	orgs, err := DBGetOrgs(db.Get(), orgID)
	if err != nil {
		responses.GetError(r, err.Error())
		return
	}

	// Return the list of organizations or a 404 if the user wasn't found.
	if orgID != -1 && len(orgs) < 1 {
		responses.GetNotFound(r)
	} else {
		responses.GetOK(r, orgs)
	}
}

// DBGetOrgs returns an organization object from the database.
func DBGetOrgs(db *sql.DB, orgID int64) (OrgList, error) {

	// If no orgID is provided (orgID is -1), all users are retreived. If
	// a orgID is given, a specific user is retreived.
	var rows *sql.Rows

	if orgID == -1 {
		queryString := `SELECT * from organization`
		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query()
		if err != nil {
			return nil, err
		}

	} else {
		queryString := `SELECT * from organization WHERE id = ?`
		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query(orgID)
		if err != nil {
			return nil, err
		}
	}

	// Read the rows into the target struct
	var objs []Org

	for rows.Next() {

		var obj Org
		err := rows.Scan(
			&obj.ID, &obj.Name)

		// Forward the error
		if err != nil {
			return nil, err
		}
		// Append object to the list
		objs = append(objs, obj)
	}

	return objs, nil
}
