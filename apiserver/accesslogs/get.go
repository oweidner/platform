//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package accesslogs

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// Get retrieves one or more user objects from the database and
// sends them back to caller.
//
func Get(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// logID is either -1 if no log ID was provided or > 0 otherwise.
	var logID int64 = -1

	// Convert the logID string to an integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		logID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.GetError(r, fmt.Sprintf("Invalid Log ID: %v", logID))
			return
		}
	}
	// Retrieve the (list of) access logs from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	logs, err := DBGetLogs(db.Get(), logID)
	if err != nil {
		responses.GetError(r, err.Error())
		return
	}

	// Return the list of logs or a 404 if the log wasn't found.
	if logID != -1 && len(logs) < 1 {
		responses.GetNotFound(r)
	} else {
		responses.GetOK(r, logs)
	}
}
