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

	"github.com/codewerft/platform/auth"
	"github.com/codewerft/platform/logging"

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
func Auth(u AuthRequest, a auth.Authenticator, r render.Render) {

	// Authenticate the user with the password provided
	user, err := a.Auth(u.Username, []byte(u.Password))
	if err != nil {
		logging.Log.Error(fmt.Sprintf("[auth] Authentication failed for user %v", u.Username))
		r.JSON(http.StatusUnauthorized,
			ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Authorization Failed"})
		return
	} else {
		logging.Log.Info(fmt.Sprintf("[auth] Authentication granted to user %v", u.Username))
	}

	// Create a new JWT token
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["org"] = user.Organization
	token.Claims["user"] = user.Username
	token.Claims["role"] = user.Role

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
