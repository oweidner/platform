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

	"github.com/codewerft/platform/apiserver/responses"

	"github.com/codewerft/platform/database"
	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// CreateUserRequest is the object that is expected by the
// CreateUser() function.
type CreateUserRequest struct {
	Firstname string `binding:"required"`
	Lastname  string `binding:"required"`
	Email     string `binding:"required"`
	Username  string `binding:"required"`
	Password  string `binding:"required"`
}

// Create creates a new User object in
// the database.
func Create(r render.Render, params martini.Params, db database.Datastore, data CreateUserRequest) {

	// Store the user object the database. In case the
	// database operation fails, an error response is sent back to the caller.
	newUser, err := DBCreateUser(db.Get(), data)
	if err != nil {
		responses.GetError(r, err.Error())
		return
	}
	// Return the user.
	responses.CreateOK(r, newUser)
}

// DBCreateUser creates a new User object in the database.
//
func DBCreateUser(db *sql.DB, data CreateUserRequest) (UserList, error) {

	stmt, err := db.Prepare(`
		INSERT account SET firstname=?, lastname=?, email=?, username=?,
        password=?`)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(
		data.Firstname, data.Lastname, data.Email, data.Username,
		data.Password)
	if err != nil {
		return nil, err
	}

	// The id of the newly generated user
	userID, _ := res.LastInsertId()
	// Retrieve the newly created object from the database and return it
	users, err := DBGetUsers(db, userID)
	if err != nil {
		return nil, err
	}
	// Return the user object.
	return users, nil
}
