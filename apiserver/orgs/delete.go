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
	"fmt"
	"net/http"
	"strconv"

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// Delete removes one or more organization objects from the database and
// sends them back to caller.
//
func Delete(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// orgID is either -1 if no organization ID was provided or > 0 otherwise.
	var orgID int64 = -1

	// Convert the organization ID string to a 64-bit integer. In case the
	// conversion fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		orgID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.DeleteError(r, fmt.Sprintf("Invalid Organization ID: %v", orgID))
			return
		}
	}
	// Use the Generic 'delete' function from the database helper package to
	// delete the user.
	// err := database.GenericDelete(db.Get(), "DELETE FROM org WHERE id = ?", orgID)
	// if err != nil {
	// 	responses.DeleteError(r, err.Error())
	// 	return
	// }

	// Respond with OK.
	responses.DeleteOK(r, "Organization deleted")
}