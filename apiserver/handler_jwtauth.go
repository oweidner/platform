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

	"github.com/dgrijalva/jwt-go"
	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

type JWTAuthAccount struct {
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type ErrorResponse struct {
	Code    int    `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
}

// JWTAuth is a custom Martini handler that checks a JWT for its validity
func JWTAuth(authEnabled interface{}, requiredRole interface{}) martini.Handler {

	// Return a handler function
	return func(context martini.Context, req *http.Request, r render.Render) {

		// If the base parameter is set to 'false', auth is disabled.
		if authEnabled.(bool) == false {
			return
		}

		// Get the caller's IP address. Used for logging.
		// originIP := req.RemoteAddr

		// Extract the JWT from the http request
		token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
			return jwtPublicKey, nil
		})

		// Token extraction didn't work. This is an error.
		if err != nil {
			r.JSON(http.StatusUnauthorized,
				ErrorResponse{
					Code:    http.StatusUnauthorized,
					Message: err.Error()})
			return
		}

		// Token is not valid. This is definetly an error.
		if token.Valid == false {
			r.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid Token. This incident will be logged."})
			return
		}

		// We have a valid token. Now we can check if user and role data
		// can be extracted, and if so, make sure the user has the correct role.
		user := token.Claims["user"]
		role := token.Claims["role"]

		if user == nil || role == nil {
			r.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid Token Data. This incident will be logged."})
			return
		}

		// Compare the roles.
		if role != requiredRole.(string) {

			// accesslogs.DBWriteLoginError(ap.Database, origin, username, dbError.Error(), userid)

			r.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid credentials for this request. This incident will be logged."})
			return
		}

		context.Map(JWTAuthAccount{Username: user.(string), Role: role.(string)})
	}
}
