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

// Get retrieves one or more plan objects from the database and
// sends them back to caller.
//
func Get(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// planID is either -1 if no plan ID was provided or > 0 otherwise.
	var planID int64 = -1

	// Convert the plan ID string to an integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		planID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.GetError(r, fmt.Sprintf("Invalid Plan ID: %v", planID))
			return
		}
	}
	// Retrieve the (list of) plans from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	plans, err := DBGetPlans(db.Get(), planID)
	if err != nil {
		responses.GetError(r, err.Error())
		return
	}

	// Return the list of plans or a 404 if the user wasn't found.
	if planID != -1 && len(plans) < 1 {
		responses.GetNotFound(r)
	} else {
		responses.GetOK(r, plans)
	}
}

// DBGetPlans returns a Plan object from the database.
func DBGetPlans(db *sql.DB, planID int64) (PlanList, error) {

	// If no userID is provided (userID is -1), all users are retreived. If
	// a planID is given, a specific user is retreived.
	var rows *sql.Rows

	if planID == -1 {
		queryString := `SELECT * from plan`
		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query()
		if err != nil {
			return nil, err
		}

	} else {
		queryString := `SELECT * from plan WHERE id = ?`
		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query(planID)
		if err != nil {
			return nil, err
		}
	}

	// Read the rows into the target struct
	var objs []Plan

	for rows.Next() {

		var obj Plan
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
