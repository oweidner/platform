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
	"errors"
	"fmt"

	"github.com/codewerft/platform/apiserver/accesslogs"
	"github.com/codewerft/platform/apiserver/account"
	"github.com/codewerft/platform/database"

	"golang.org/x/crypto/bcrypt"
)

// The DefaultAuthProvider provides an authenticator that checks an auth
// request agains a static string.
type DefaultAuthProvider struct {
	Name        string
	RootAccount account.Account
	Database    database.Datastore
}

// NewDefaultAuthProvider creates a new StaticAuthProvider object.
func NewDefaultAuthProvider(ds database.Datastore, rootAccount account.Account) *DefaultAuthProvider {
	return &DefaultAuthProvider{Name: "MySQLAuthProvider", Database: ds, RootAccount: rootAccount}
}

// Auth implements the AuthProvider interface.
func (ap *DefaultAuthProvider) Auth(origin string, accountname string, password []byte) error {

	var finalError error

	// Search the database for a matching entry
	db := ap.Database.Get()
	var accountUsername, accountPassword string
	sqlErr := db.QueryRow(`
		SELECT username, password FROM platform_account
        WHERE username = ?`, accountname).Scan(&accountUsername, &accountPassword)
	if sqlErr == nil {
		bcryptErr := bcrypt.CompareHashAndPassword([]byte(accountPassword), password)
		if bcryptErr == nil {
			// Account authenticated successfully.
			accesslogs.DBWriteLoginOK(ap.Database.Get(), origin, accountname)
			return nil
		}
		finalError = bcryptErr

	} else {
		finalError = sqlErr
	}

	// Account wasn't found in the database. Let's check if it's the root account.
	if accountUsername == "" {
		if accountname == ap.RootAccount.Username {
			bcryptErr := bcrypt.CompareHashAndPassword([]byte(ap.RootAccount.Password), password)
			if bcryptErr == nil {
				// Root account authenticated successfully.
				accesslogs.DBWriteLoginOK(ap.Database.Get(), origin, accountname)
				return nil
			}
			finalError = bcryptErr
		} else {
			finalError = fmt.Errorf("Unknown username '%v'", accountname)
		}
	}

	// Neither datbase nor root account authentication was successfull.
	accesslogs.DBWriteLoginError(ap.Database.Get(), origin, accountname, finalError.Error())
	return errors.New("Authentication failed.")
}
