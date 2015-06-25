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

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// CreatePlanRequest is the object that is expected by the
// Create() function.
type CreatePlanRequest struct {
	Name string `binding:"required"`
}

// Create creates a new Account object in
// the database.
func Create(r render.Render, params martini.Params, db database.Datastore, data CreatePlanRequest) {

	// Store the plan object in the database. In case the
	// database operation fails, an error response is sent back to the caller.
	newPlan, err := DBCreatePlan(db.Get(), data)
	if err != nil {
		responses.CreateError(r, err.Error())
		return
	}
	// Return the account.
	responses.CreateOK(r, newPlan)
}

// DBCreatePlan creates a new Account object in the database.
//
func DBCreatePlan(db *sql.DB, data CreatePlanRequest) (PlanList, error) {

	stmt, err := db.Prepare(`INSERT platform_plan SET name=?`)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(data.Name)
	if err != nil {
		return nil, err
	}

	// The id of the newly generated plan
	planID, _ := res.LastInsertId()
	// Retrieve the newly created object from the database and return it
	plans, err := DBGetPlans(db, planID)
	if err != nil {
		return nil, err
	}
	// Return the account object.
	return plans, nil
}
