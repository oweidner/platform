//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package orgplans

import (
	"net/http"

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/apiserver/utils"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// List returns the list of available resources.
//
func List(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	db.GetDBMap().AddTableWithName(OrganisationPlanAssoc{}, "platform_organisation_plan_assoc").SetKeys(true, "id")

	// Parse the resource ID into an int64
	resourceID, parseError := utils.ParseResourceID(params["p1"])
	if parseError != nil {
		responses.Error(r, parseError.Error())
	}

	// Retreive the requested resource from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	var plans OrganisationPlanAssocList
	_, err := db.GetDBMap().Select(&plans, "SELECT * FROM platform_organisation_plan_assoc WHERE organisation_id=?", resourceID)
	if err != nil {
		responses.Error(r, err.Error())
		return
	}
	responses.OKStatusPlusData(r, plans, len(plans))
}

// Get retrieves one or more resource objects from the database and
// sends them back to caller.
//
func Get(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	db.GetDBMap().AddTableWithName(OrganisationPlanAssoc{}, "platform_organisation_plan_assoc").SetKeys(true, "id")

	// Parse the resource ID into an int64
	resourceID, parseError := utils.ParseResourceID(params["p1"])
	if parseError != nil {
		responses.Error(r, parseError.Error())
	}

	// Parse the resource ID into an int64
	planID, parseError := utils.ParseResourceID(params["p2"])
	if parseError != nil {
		responses.Error(r, parseError.Error())
	}

	// Retreive the list of all resources from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	var plan OrganisationPlanAssoc
	err := db.GetDBMap().SelectOne(&plan, "SELECT * FROM platform_organisation_plan_assoc WHERE organisation_id=? AND id=?", resourceID, planID)
	if err != nil {
		responses.Error(r, err.Error())
		return
	}
	responses.OKStatusPlusData(r, plan, 1)
}
