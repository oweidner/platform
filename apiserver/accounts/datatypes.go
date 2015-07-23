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

import "database/sql"

// Account represents an Account object as it exists
// in the database.
type Account struct {
	ID                int64          `db:"id"`
	Disabled          int            `db:"disabled"`
	Username          string         `db:"username"`
	Password          string         `db:"password"`
	Firstname         string         `db:"firstname"`
	Lastname          string         `db:"lastname"`
	ContactEmail      string         `db:"contact_email"`
	ContactPhone      sql.NullString `db:"contact_phone"`
	AddressStreet1    sql.NullString `db:"address_street_1"`
	AddressStreet2    sql.NullString `db:"address_street_2"`
	AddressZIP        sql.NullString `db:"address_zip"`
	AddressCity       sql.NullString `db:"address_city"`
	AddressCountry    sql.NullString `db:"address_country"`
	BankAccountHolder sql.NullString `db:"bank_account_holder"`
	BankIBAN          sql.NullString `db:"bank_IBAN"`
	BankBIC           sql.NullString `db:"bank_BIC"`
	BankName          sql.NullString `db:"bank_name"`
	Roles             []string       `db:"-"`
}

// AccountList represents a list of Account objects.
type AccountList []Account
