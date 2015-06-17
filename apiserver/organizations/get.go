//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.
//

package organizations

import (
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

	// Convert the orgID string to an integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		orgID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.GetError(r, "Invalid Organization ID.")
			return
		}
	}
	// Retrieve the (list of) organizations from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	orgs, err := DBGetOrganizations(db.Get(), orgID)
	if err != nil {
		responses.GetError(r, err.Error())
		return
	}

	// Return the user.
	responses.GetOK(r, orgs)
}
