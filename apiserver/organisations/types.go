//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package organisations

import "gopkg.in/guregu/null.v2"

// Organisation represents an organization object as it exists in the database.
type Organisation struct {
	ID             int64       `db:"id"`
	Orgname        string      `db:"orgname"`
	Name           string      `db:"name"`
	ContactEmail   string      `db:"contact_email"`
	ContactPhone   null.String `db:"contact_phone"`
	AddressStreet1 null.String `db:"address_street_1"`
	AddressStreet2 null.String `db:"address_street_2"`
	AddressCity    null.String `db:"address_city"`
	AddressZIP     null.String `db:"address_zip"`
	AddressCountry null.String `db:"address_country"`
}

// OrganisationList represents a list of Organization object.
type OrganisationList []Organisation
