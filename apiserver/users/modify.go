//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package users

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// ModifyUserRequest is the object that is expected by the
// Modify() function.
type ModifyUserRequest struct {
	Firstname string `binding:"required"`
	Lastname  string `binding:"required"`
	Email     string `binding:"required"`
	Password  string `binding:"required"`
}

// Modify modifies a user object in the database.
//
func Modify(r render.Render, params martini.Params, db database.Datastore, data ModifyUserRequest) {

	// userID is either -1 if no user ID was provided or > 0 otherwise.
	var userID int64 = -1

	// Convert the userid string to an integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		userID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.ModifyError(r, fmt.Sprintf("Invalid User ID: %v", userID))
			return
		}
	}
	// Retrieve the (list of) users from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	modifiedUser, err := DBModifyUser(db.Get(), userID, data)
	if err != nil {
		responses.GetError(r, err.Error())
		return
	}

	// Return the user.
	responses.ModifyOK(r, modifiedUser)
}

// DBModifyUser modifies a User object in the database.
//
func DBModifyUser(db *sql.DB, userID int64, data ModifyUserRequest) (UserList, error) {

	stmt, err := db.Prepare(`
		UPDATE user SET firstname=?, lastname=?, email=?, password=?
        WHERE id=?`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(
		data.Firstname, data.Lastname, data.Email, data.Password, userID)
	if err != nil {
		return nil, err
	}

	// Retrieve the modified object from the database and return it
	users, err := DBGetUsers(db, userID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
