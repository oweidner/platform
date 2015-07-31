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
	"net/http"
	"time"

	"gopkg.in/guregu/null.v2"

	"github.com/codewerft/platform/apiserver/accounts"
	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/auth"
	"github.com/codewerft/platform/database"
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
	data, err := a.Auth(originIP, u.Username, []byte(u.Password))
	if err != nil {
		r.JSON(http.StatusUnauthorized,
			ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Authorization Failed"})
		return
	}

	// Create a new JWT token
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["userid"] = data.ID     //account.Username
	token.Claims["user"] = data.Username //account.Username
	token.Claims["role"] = "data.Role"   //account.Role

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

	// Extract the userid.
	user := token.Claims["user"]
	userid := token.Claims["userid"]

	// account holds the data returned to the caller
	var account accounts.Account

	// If the UserID is 0, this is the platform admin user.
	if userid == 0 {
		// Create a 'fake' admin account
		account = accounts.Account{
			ID:           1,
			Firstname:    null.StringFrom("Platform"),
			Lastname:     null.StringFrom("Superuser"),
			ContactEmail: null.StringFrom("platform.codewerft.net"),
			Username:     null.StringFrom(user.(string)),
			Roles:        []string{"admin"}}
	} else {
		// Query the database for the given UserID
		err := db.GetDBMap().SelectOne(&account, "SELECT * FROM platform_account WHERE _deleted=0 AND id=?", userid)
		if err != nil {
			responses.Error(r, err.Error())
			r.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error()})
		}
	}

	r.JSON(http.StatusOK, account)
}
