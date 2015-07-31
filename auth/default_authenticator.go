//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package auth

import (
	"github.com/codewerft/platform/apiserver/accesslogs"
	"github.com/codewerft/platform/apiserver/accounts"
	"github.com/codewerft/platform/database"
	"gopkg.in/guregu/null.v2"

	"golang.org/x/crypto/bcrypt"
)

// The DefaultAuthProvider provides an authenticator that checks an auth
// request agains a static string.
type DefaultAuthProvider struct {
	Name        string
	RootAccount accounts.Account
	Database    database.Datastore
}

// NewDefaultAuthProvider creates a new StaticAuthProvider object.
func NewDefaultAuthProvider(ds database.Datastore, rootAccount accounts.Account) *DefaultAuthProvider {
	return &DefaultAuthProvider{Name: "MySQLAuthProvider", Database: ds, RootAccount: rootAccount}
}

// Auth implements the AuthProvider interface.
// It authenticates a given ussername and password combination against:
// - the admin user(s) defined in the configuration file
// - the the account database.
func (ap *DefaultAuthProvider) Auth(origin string, username string, password []byte) (accounts.Account, error) {

	// Retreive the list of all resources from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	var account accounts.Account
	var userid null.Int

	dbError := ap.Database.GetDBMap().SelectOne(&account, "SELECT * FROM platform_account WHERE _deleted=0 AND username=?", username)
	if dbError != nil {
		accesslogs.DBWriteLoginError(ap.Database, origin, username, dbError.Error(), userid)
		return accounts.Account{}, dbError
	}

	// user was found in he database, hence there's a userid
	userid = null.IntFrom(account.ID)

	// Compare the hashed passwords
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(account.Password), password)
	if bcryptErr != nil {
		accesslogs.DBWriteLoginError(ap.Database, origin, username, bcryptErr.Error(), userid)
		return accounts.Account{}, bcryptErr
	}

	// Account authenticated successfully.
	accesslogs.DBWriteLoginOK(ap.Database, origin, username, userid)
	return account, nil
}
