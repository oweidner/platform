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

import "gopkg.in/guregu/null.v2"

// OrganisationPlanAssoc represents the plans associated to an organisation
type OrganisationPlanAssoc struct {
	ID              int64       `db:"id"`
	OrganisationID  int64       `db:"organisation_id"`
	PlanID          int64       `db:"plan_id"`
	PlanName        null.String `db:"plan_name"`
	StartDate       null.String `db:"start_date"`
	EndDate         null.String `db:"end_date"`
	TerminationDate null.String `db:"termination_date"`

	Terminated null.Bool `db:"terminated"`
}

// OrganisationPlanAssoc represents a list of OrganisationPlanAssoc object.
type OrganisationPlanAssocList []OrganisationPlanAssoc
