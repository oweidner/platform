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

import "gopkg.in/guregu/null.v2"

// Plan represents a Plan object as it exists in the database.
type Plan struct {
	ID              int64       `db:"id"`
	Deleted         null.Bool   `db:"_deleted"`
	Name            string      `db:"name"`
	Description     null.String `db:"description"`
	Parameters      null.String `db:"parameters"`
	Rate            float64     `db:"rate"`
	VATPercentage   null.Float  `db:"vat_percentage"`
	BillingInterval int         `db:"billing_interval"`
	MinDuration     int         `db:"min_duration"`
}

// PlanList represents a list of Organization object.
type PlanList []Plan
