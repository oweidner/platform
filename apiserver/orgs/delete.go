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

// Delete removes one or more Org objects from the database and
// sends them back to caller.
//
func Delete(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// OrgID is either -1 if no Org ID was provided or > 0 otherwise.
	var OrgID int64 = -1

	// Convert the Org ID string to a 64-bit integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		OrgID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.DeleteError(r, fmt.Sprintf("Invalid Organisation ID: %v", OrgID))
			return
		}
	}

	// Delete the Org object from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	err := DBDeleteOrg(db.Get(), OrgID)
	if err != nil {
		responses.DeleteError(r, err.Error())
		return
	}

	// Return the modified Org.
	responses.DeleteOK(r, "Organisation removed from Database")
}

// DBDeleteOrg removes the Org from the MySQL database.
func DBDeleteOrg(db *sql.DB, OrgID int64) error {

	stmt, err := db.Prepare(`
		DELETE from platform_organisation WHERE id=?`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(OrgID)
	if err != nil {
		return err
	}

	return nil
}
