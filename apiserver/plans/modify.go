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
	"strconv"

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// ModifyPlanRequest is the object that is expected by the
// Modify() function.
type ModifyPlanRequest struct {
	Name            string
	Description     string  `json:",omitempty"`
	Parameters      string  `json:",omitempty"`
	Rate            float64 `json:",string"`
	VATPercentage   float64 `json:",string,omitempty"`
	BillingInterval uint64  `json:",string"`
	MinDuration     uint64  `json:",string"`
}

// Modify modifies a plan object in the database.
//
func Modify(r render.Render, params martini.Params, db database.Datastore, data ModifyPlanRequest) {

	// planID is either -1 if no planm ID was provided or > 0 otherwise.
	var planID int64 = -1

	// Convert the plan ID string to a 64-bit integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		planID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.ModifyError(r, fmt.Sprintf("Invalid Plan ID: %v", planID))
			return
		}
	}
	// Update the plan object in the database. In case the
	// database operation fails, an error response is sent back to the caller.
	modifiedPlan, err := DBModifyPlan(db.Get(), planID, data)
	if err != nil {
		responses.ModifyError(r, err.Error())
		return
	}

	// Return the modified plan.
	responses.ModifyOK(r, modifiedPlan)
}

// DBModifyPlan modifies a Account object in the database.
//
func DBModifyPlan(db *sql.DB, planID int64, data ModifyPlanRequest) (PlanList, error) {

	stmt, err := db.Prepare(`
		UPDATE platform_plan SET name=?, description=?, rate=?, billing_interval=?, min_duration=?, vat_percentage=?, parameters=? WHERE id=?`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(data.Name, data.Description, data.Rate, data.BillingInterval, data.MinDuration, data.VATPercentage, data.Parameters, planID)
	if err != nil {
		return nil, err
	}

	// Retrieve the modified object from the database and return it
	account, err := DBGetPlans(db, planID)
	if err != nil {
		return nil, err
	}

	return account, nil
}
