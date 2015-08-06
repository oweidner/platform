//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package authentication

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

type JWTConfig struct {
	PublicKey  []byte
	PrivateKey []byte
	Expiration int
}

type JWTConfigIF interface {
	Get() JWTConfig
}

type TheJWTConfig struct {
	Config JWTConfig
}

func (c TheJWTConfig) Get() JWTConfig {
	// ds.Session.Close()
	return c.Config
}

type UserInfo struct {
	UserID           int64   `db:"userid"`
	Username         string  `db:"username"`
	Firstname        string  `db:"firstname"`
	Lastname         string  `db:"lastname"`
	OrganisationID   int64   `db:"org_id"`
	OrganisationName string  `db:"org_name"`
	RoleID           int64   `db:"role_id"`
	RoleName         string  `db:"role_name"`
	Expires          float64 `db:"-"`
	IsAdmin          bool    `db:"-"`
	Password         string  `db:"password"`
}

type ErrorResponse struct {
	Code    int    `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
}

// JWTAuth is a custom Martini handler that checks a JWT for its validity
func JWTAuth(jwtcfg JWTConfig, requiredRole interface{}) martini.Handler {

	// Return a handler function
	return func(context martini.Context, req *http.Request, r render.Render) {

		// Extract the JWT from the http request
		token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
			return jwtcfg.PublicKey, nil
		})

		if err != nil {
			// Token extraction didn't work. This is an error.
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
		if token.Claims["user_id"] == nil || token.Claims["user"] == nil ||
			token.Claims["firstname"] == nil || token.Claims["lastname"] == nil ||
			token.Claims["org_id"] == nil || token.Claims["org_name"] == nil ||
			token.Claims["role_id"] == nil || token.Claims["role_name"] == nil ||
			token.Claims["exp"] == nil || token.Claims["is_admin"] == nil {

			r.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid Token Data - missing fields. This incident will be logged."})
			return
		}

		// Compare the roles.
		if requiredRole != nil && token.Claims["org_role"] != requiredRole.(string) {

			r.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid role for this request. This incident will be logged."})
			return
		}

		// Explicit conversion due to weird JWT.Claims behavior.
		userID, _ := strconv.ParseInt(token.Claims["user_id"].(string), 10, 64)
		roleID, _ := strconv.ParseInt(token.Claims["role_id"].(string), 10, 64)
		organisationID, _ := strconv.ParseInt(token.Claims["org_id"].(string), 10, 64)

		context.Map(UserInfo{
			UserID:           userID,
			Username:         token.Claims["user"].(string),
			Firstname:        token.Claims["firstname"].(string),
			Lastname:         token.Claims["lastname"].(string),
			OrganisationID:   organisationID,
			OrganisationName: token.Claims["org_name"].(string),
			RoleID:           roleID,
			RoleName:         token.Claims["role_name"].(string),
			Expires:          token.Claims["exp"].(float64),
			IsAdmin:          token.Claims["is_admin"].(bool)})
	}
}
