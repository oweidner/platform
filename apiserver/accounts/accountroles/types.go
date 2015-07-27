//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package accountroles

// SQLTableName holds the table name for the account roles
var SQLTableName = "platform_account_organisation_role"

// AccountOrganisationRole represents the role of an account within an
// Organisation.
type AccountOrganisationRole struct {
	ID               int64  `db:"id"`
	AccountID        int64  `db:"account_id"`
	OrganisationID   int64  `db:"organisation_id"`
	RoleID           int64  `db:"role_id"`
	OrganisationName string `db:"-"`
	RoleName         string `db:"-"`
}

// AccountOrganisationRoleList is just a shortcut
type AccountOrganisationRoleList []AccountOrganisationRole
