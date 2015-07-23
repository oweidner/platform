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

// Organisation represents an organization object as it exists in the database.
type Organisation struct {
	ID             int64  `db:"id"`
	Orgname        string `db:"orgname"`
	Name           string `db:"name"`
	ContactEmail   string `db:"contact_email"`
	ContactPhone   string `db:"contact_phone"`
	AddressStreet1 string `db:"address_street_1"`
	AddressStreet2 string `db:"address_street_2"`
	AddressCity    string `db:"address_city"`
	AddressZIP     string `db:"address_zip"`
	AddressCountry string `db:"address_country"`
}

// OrganisationList represents a list of Organization object.
type OrganisationList []Organisation
