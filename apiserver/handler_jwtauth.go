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
func JWTAuth(base interface{}) martini.Handler {

	return func(context martini.Context, req *http.Request, r render.Render) {
		// If the base parameter is set to 'false', auth is disabled.
		if base.(bool) == false {
			return
		}

		// Extract the JWT from the http request
		token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
			return jwtPublicKey, nil
		})

		if err != nil {
			r.JSON(http.StatusUnauthorized,
				ErrorResponse{
					Code:    http.StatusUnauthorized,
					Message: err.Error()})
			return
		}

		if token.Valid == false {
			r.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid Token"})
			return
		}

		user := token.Claims["user"]
		role := token.Claims["role"]

		if user == nil {
			r.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid Token Data (user)"})
			return
		}
		context.Map(JWTAuthAccount{Username: user.(string), Role: role.(string)})
	}
}
