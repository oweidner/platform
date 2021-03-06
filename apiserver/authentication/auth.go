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

	"github.com/oweidner/platform/apiserver/accesslogs"
	"github.com/oweidner/platform/apiserver/responses"
	"github.com/oweidner/platform/database"

	"gopkg.in/dgrijalva/jwt-go.v2"
	"github.com/martini-contrib/render"
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
		SELECT a.id AS userid, a.username, a.firstname, a.lastname, a.password, aor.organisation_id AS org_id, org.name AS org_name, aor.role_id AS role_id, r.name AS role_name
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
	token.Claims["firstname"] = userInfo[0].Firstname
	token.Claims["lastname"] = userInfo[0].Lastname
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
