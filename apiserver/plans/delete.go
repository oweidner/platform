//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package plans

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

// Delete removes one or more plan objects from the database and
// sends them back to caller.
//
func Delete(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// planID is either -1 if no plan ID was provided or > 0 otherwise.
	var planID int64 = -1

	// Convert the plan ID string to a 64-bit integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		planID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.DeleteError(r, fmt.Sprintf("Invalid Plan ID: %v", planID))
			return
		}
	}

	// Delete the plan object from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	err := DBDeletePlan(db.Get(), planID)
	if err != nil {
		responses.DeleteError(r, err.Error())
		return
	}

	// Return the modified plan.
	responses.DeleteOK(r, "Plan deleted")
}

// DBDeletePlan removes the plan from the MySQL database.
func DBDeletePlan(db *sql.DB, planID int64) error {

	stmt, err := db.Prepare(`
		DELETE from platform_plan WHERE id=?`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(planID)
	if err != nil {
		return err
	}

	return nil
}
