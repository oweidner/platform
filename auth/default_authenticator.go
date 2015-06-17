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

	"golang.org/x/crypto/bcrypt"
)

// The DefaultAuthProvider provides an authenticator that checks an auth
// request agains a static string.
type DefaultAuthProvider struct {
	Organization string
	Name         string
	UserList     map[string]User
	Closed       bool
}

// NewDefaultAuthProvider creates a new StaticAuthProvider object.
func NewDefaultAuthProvider(userList map[string]User) *DefaultAuthProvider {
	return &DefaultAuthProvider{Name: "StaticAuthProvider", UserList: userList, Closed: false}
}

// Auth implements the AuthProvider interface.
func (ap *DefaultAuthProvider) Auth(organization string, username string, password []byte) (u User, e error) {
	if ap.Closed == true {
		return User{}, errors.New("auth provider closed")
	}

	if val, ok := ap.UserList[username]; ok {
		err := bcrypt.CompareHashAndPassword(val.Password, password)
		if err == nil {
			return val, nil
		}
	}
	return User{}, errors.New("wrong password")
}

// Close implements the AuthProvider interface.
func (ap *DefaultAuthProvider) Close() {
	ap.Closed = true
}
