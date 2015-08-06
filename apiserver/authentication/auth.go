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
	"time"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/guregu/null.v2"

	"github.com/codewerft/platform/apiserver/accesslogs"
	"github.com/codewerft/platform/apiserver/accounts"
	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"
	"github.com/go-martini/martini"

	"github.com/dgrijalva/jwt-go"
	"github.com/gavv/martini-render"
)

// AuthRequest is the object that is sent to us in order to request
// a new authentication token.
type AuthRequest struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

// AuthResponse is the object that is sent by us as a response to a
// successful authentication request.
type AuthResponse struct {
	Token string `json:"token" binding:"required"`
}

// Auth is called for every POST request on the /auth resource.
func Auth(u AuthRequest, r render.Render, req *http.Request, db database.Datastore, jwtcfg JWTConfigIF) {

	// Get the caller's IP address
	originIP := req.RemoteAddr

	// Try to retreive the user and his affiliation(s) from the database.
	userInfo := []UserInfo{}

	_, dbError := db.GetDBMap().Select(&userInfo, `
		SELECT a.id AS userid, a.username, a.password, aor.organisation_id AS org_id, org.name AS org_name, aor.role_id AS role_id, r.name AS role_name
		  FROM platform_account a LEFT JOIN platform_account_organisation_role aor ON (a.id = aor.account_id),
			platform_organisation org, platform_role r
			WHERE org.id = aor.organisation_id AND r.id = aor.role_id AND username=?;`, u.Username)

	// SQL error - return an Unauthorized error.
	if dbError != nil {
		responses.AuthenticationError(r, "Authorization failed", dbError.Error())
		return
	}

	// Result set length empty - return an Unauthorized error.
	if len(userInfo) < 1 {
		responses.AuthenticationError(r, "Authorization failed", "SQL: Empty result set.")
		return
	}

	// Result set length greater 1 - multiple organisations currently not supported. return an Unauthorized error.
	if len(userInfo) > 1 {
		responses.AuthenticationError(r, "Authorization failed", "SQL: User is associated with multiple organisations. This is not supported yet")
		return
	}

	// Found an account by the given name. Now compare the hashed passwords
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(userInfo[0].Password), []byte(u.Password))
	if bcryptErr != nil {
		accesslogs.DBWriteLoginError(db, originIP, userInfo[0].Username, bcryptErr.Error(), null.IntFrom(userInfo[0].UserID))
		responses.AuthenticationError(r, "Authorization failed", bcryptErr.Error())
		return
	}

	// Create a new JWT token. The JWT token contains all information
	// required to authenticate a token-holder on a per user / per resource /
	// per role base.
	token := jwt.New(jwt.GetSigningMethod("RS256"))

	token.Claims["user_id"] = strconv.FormatInt(userInfo[0].UserID, 10)
	token.Claims["user"] = userInfo[0].Username
	token.Claims["org_id"] = strconv.FormatInt(userInfo[0].OrganisationID, 10)
	token.Claims["org_name"] = userInfo[0].OrganisationName
	token.Claims["role_id"] = strconv.FormatInt(userInfo[0].RoleID, 10)
	token.Claims["role_name"] = userInfo[0].RoleName
	token.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(jwtcfg.Get().Expiration)).Unix()

	if userInfo[0].UserID == -1 {
		token.Claims["is_admin"] = true
	} else {
		token.Claims["is_admin"] = false
	}

	tokenString, err := token.SignedString(jwtcfg.Get().PrivateKey)
	if err != nil {
		r.HTML(201, "error", nil)
		return
	}
	// Return the token
	r.JSON(http.StatusOK, AuthResponse{Token: tokenString})
}

// GetSelf retrieves the account referenced in the auth token.
//
func GetSelf(req *http.Request, params martini.Params, r render.Render, db database.Datastore, jwtcfg JWTConfigIF) {

	// Extract the JWT from the http request
	token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		return jwtcfg.Get().PublicKey, nil
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
	userid := token.Claims["userid"].(string)
	// role := token.Claims["role"]

	// account holds the data returned to the caller
	var account accounts.Account

	// If the UserID is 0, this is the platform admin user.
	if userid == "-1" {
		// Create a 'fake' admin account
		account = accounts.Account{
			ID:           -1,
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
			return
		}
	}

	r.JSON(http.StatusOK, account)
}
