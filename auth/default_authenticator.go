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
	"log"

	"github.com/codewerft/platform/apiserver/accesslogs"
	"github.com/codewerft/platform/apiserver/users"
	"github.com/codewerft/platform/database"

	"golang.org/x/crypto/bcrypt"
)

// The DefaultAuthProvider provides an authenticator that checks an auth
// request agains a static string.
type DefaultAuthProvider struct {
	Name     string
	RootUser users.User
	Database database.Datastore
}

// NewDefaultAuthProvider creates a new StaticAuthProvider object.
func NewDefaultAuthProvider(ds database.Datastore, rootUser users.User) *DefaultAuthProvider {
	return &DefaultAuthProvider{Name: "MySQLAuthProvider", Database: ds, RootUser: rootUser}
}

// Auth implements the AuthProvider interface.
func (ap *DefaultAuthProvider) Auth(origin string, username string, password []byte) error {
	// Search the database for a matching entry
	db := ap.Database.Get()
	var accountUsername, accountPassword string

	sqlErr := db.QueryRow(`
		SELECT username, password FROM platform_account
        WHERE username = ?`, username).Scan(&accountUsername, &accountPassword)
	if sqlErr == nil {
		bcryptErr := bcrypt.CompareHashAndPassword([]byte(accountPassword), password)
		if bcryptErr != nil {
			// TODO: Handle the BCrypt error better
			log.Print(bcryptErr)
			return errors.New("Authentication failed (bcrypt error)")
		}
		// User authenticated successfully.
		accesslogs.DBWriteLoginOK(ap.Database.Get(), origin, username)
		return nil
	} else {
		// TODO: handle better
		log.Print(sqlErr)
	}

	// User wasn't found in the database. Let's check if it's the
	// root user.
	if username == ap.RootUser.Username {
		bcryptErr := bcrypt.CompareHashAndPassword([]byte(ap.RootUser.Password), password)
		if bcryptErr != nil {
			// TODO: Handle the BCrypt error better
			log.Fatal(bcryptErr)
			return errors.New("Authentication failed (bcrypt error)")
		}
		// User authenticated successfully.
		accesslogs.DBWriteLoginOK(ap.Database.Get(), origin, username)
		return nil
	}

	// Neither datbase nor root user authentication was successfull.
	return errors.New("Authentication failed.")
}
