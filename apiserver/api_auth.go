//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package apiserver

import (
	"fmt"
	"net/http"
	"time"

	"gopkg.in/guregu/null.v2"

	"github.com/codewerft/platform/apiserver/accounts"
	"github.com/codewerft/platform/auth"
	"github.com/codewerft/platform/database"
	"github.com/codewerft/platform/logging"
	"github.com/go-martini/martini"

	"github.com/dgrijalva/jwt-go"
	"github.com/gavv/martini-render"
)

// AuthRequest is the object that is sent to us in order to request
// a new authentication token.
type AuthRequest struct {
	Username string `json:"username"     binding:"required"`
	Password string `json:"password"     binding:"required"`
}

// AuthResponse is the object that is sent by us as a response to a
// successful authentication request.
type AuthResponse struct {
	Token string `json:"token" binding:"required"`
}

// Auth is called for every POST request on the /auth resource.
func Auth(u AuthRequest, a auth.Authenticator, r render.Render, req *http.Request) {

	// Get the caller's IP address
	originIP := req.RemoteAddr

	// Authenticate the account with the password provided
	err := a.Auth(originIP, u.Username, []byte(u.Password))
	if err != nil {
		logging.Log.Error(fmt.Sprintf("[auth] Authentication failed for account %v", u.Username))
		r.JSON(http.StatusUnauthorized,
			ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Authorization Failed"})
		return
	}

	logging.Log.Info(fmt.Sprintf("[auth] Authentication granted to account %v", u.Username))

	// Create a new JWT token
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["org"] = "ORG"       //account.Organization
	token.Claims["user"] = u.Username //account.Username
	token.Claims["role"] = "ROLE"     //account.Role

	// Expire in 60 mins
	token.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(jwtExpiration)).Unix()

	tokenString, err := token.SignedString(jwtPrivateKey)
	if err != nil {
		r.HTML(201, "error", nil)
		return
	}
	// Return the token
	r.JSON(http.StatusOK, AuthResponse{Token: tokenString})
}

// GetSelf retrieves the account referenced in the auth token.
//
func GetSelf(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// Extract the JWT from the http request
	token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		return jwtPublicKey, nil
	})

	// Checks if a token was extracted successfully; returns an error if not.
	if err != nil {
		r.JSON(http.StatusUnauthorized,
			ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error()})
		return
	}

	// Checks is the token is valid; returns an error if not.
	if token.Valid == false {
		r.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Invalid Token"})
		return
	}

	// Extract the username.
	user := token.Claims["user"]

	// Create a 'fake' admin account
	obj := accounts.Account{
		ID:           1,
		Firstname:    null.StringFrom("Platform"),
		Lastname:     null.StringFrom("Superuser"),
		ContactEmail: null.StringFrom("platform.codewerft.net"),
		Username:     null.StringFrom(user.(string)),
		Roles:        []string{"admin"}}

	type Account struct {
		ID           int64
		Firstname    string
		Lastname     string
		ContactEmail string
		Username     string
		Password     string
		Roles        []string
	}

	r.JSON(http.StatusOK, obj)
}
