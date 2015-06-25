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

	"github.com/codewerft/platform/apiserver/accesslogs"
	"github.com/codewerft/platform/database"

	"golang.org/x/crypto/bcrypt"
)

// The DefaultAuthProvider provides an authenticator that checks an auth
// request agains a static string.
type DefaultAuthProvider struct {
	Name     string
	Database database.Datastore
	UserList map[string]User
	Closed   bool
}

// NewDefaultAuthProvider creates a new StaticAuthProvider object.
func NewDefaultAuthProvider(userList map[string]User, ds database.Datastore) *DefaultAuthProvider {
	return &DefaultAuthProvider{Name: "StaticAuthProvider", UserList: userList, Database: ds, Closed: false}
}

// Auth implements the AuthProvider interface.
func (ap *DefaultAuthProvider) Auth(origin string, username string, password []byte) (u User, e error) {
	if ap.Closed == true {
		return User{}, errors.New("auth provider closed")
	}

	if val, ok := ap.UserList[username]; ok {
		err := bcrypt.CompareHashAndPassword(val.Password, password)
		if err == nil {
			msg := "Authenticated"

			entry := accesslogs.CreateAccessLogEntryRequest{
				Origin: origin, Level: "INFO",
				Event: msg, AccountID: 1}
			accesslogs.DBCreateAccessLogEntry(ap.Database.Get(), entry)

			return val, nil
		}
	}
	// Write the authentication error to the access logs.
	msg := "Authentication failed: wrong password"

	entry := accesslogs.CreateAccessLogEntryRequest{
		Origin: origin, Level: "ERROR",
		Event: msg, AccountID: 1}
	accesslogs.DBCreateAccessLogEntry(ap.Database.Get(), entry)

	// Return the same error to the caller.
	return User{}, errors.New(msg)
}

// Close implements the AuthProvider interface.
func (ap *DefaultAuthProvider) Close() {
	ap.Closed = true
}
