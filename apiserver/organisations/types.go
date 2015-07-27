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
	ID                 int64       `db:"id"`
	Deleted            null.Bool   `db:"_deleted"`
	Disabled           null.Bool   `db:"disabled"`
	Orgname            null.String `db:"orgname"`
	Name               null.String `db:"name"`
	ContactEmail       null.String `db:"contact_email"`
	ContactPerson      null.String `db:"contact_person"`
	ContactPersonTitle null.String `db:"contact_person_title"`
	ContactPhone       null.String `db:"contact_phone"`
	AddressStreet1     null.String `db:"address_street_1"`
	AddressStreet2     null.String `db:"address_street_2"`
	AddressCity        null.String `db:"address_city"`
	AddressZIP         null.String `db:"address_zip"`
	AddressCountry     null.String `db:"address_country"`
	BankAccountHolder  null.String `db:"bank_account_holder"`
	BankIBAN           null.String `db:"bank_iban"`
	BankBIC            null.String `db:"bank_bic"`
	BankName           null.String `db:"bank_name"`
}

// OrganisationList represents a list of Organization object.
type OrganisationList []Organisation
