//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package accounts

import "gopkg.in/guregu/null.v2"

// Account represents an Account object as it exists
// in the database.
type Account struct {
	ID                int64       `db:"id"`
	Deleted           null.Bool   `db:"_deleted"`
	Disabled          null.Bool   `db:"disabled"`
	Username          null.String `db:"username"`
	Password          string      `db:"password"`
	Firstname         null.String `db:"firstname"`
	Lastname          null.String `db:"lastname"`
	Title             null.String `db:"title"`
	ContactEmail      null.String `db:"contact_email"`
	ContactPhone      null.String `db:"contact_phone"`
	AddressStreet1    null.String `db:"address_street_1"`
	AddressStreet2    null.String `db:"address_street_2"`
	AddressZIP        null.String `db:"address_zip"`
	AddressCity       null.String `db:"address_city"`
	AddressCountry    null.String `db:"address_country"`
	BankAccountHolder null.String `db:"bank_account_holder"`
	BankIBAN          null.String `db:"bank_iban"`
	BankBIC           null.String `db:"bank_bic"`
	BankName          null.String `db:"bank_name"`
	Roles             []string    `db:"-"`
}

// AccountList represents a list of Account objects.
type AccountList []Account
