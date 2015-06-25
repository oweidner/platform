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
	"net/http"
	"strconv"

	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// Get retrieves one or more user objects from the database and
// sends them back to caller.
//
func Get(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// userID is either -1 if no user ID was provided or > 0 otherwise.
	var userID int64 = -1

	// Convert the userid string to an integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		userID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.GetError(r, fmt.Sprintf("Invalid User ID: %v", userID))
			return
		}
	}
	// Retrieve the (list of) users from the database. In case the
	// database operation fails, an error response is sent back to the caller.
	users, err := DBGetUsers(db.Get(), userID)
	if err != nil {
		responses.GetError(r, err.Error())
		return
	}

	// Return the list of users or a 404 if the user wasn't found.
	if userID != -1 && len(users) < 1 {
		responses.GetNotFound(r)
	} else {
		responses.GetOK(r, users)
	}
}

// DBGetUsers returns a User object from the database.
func DBGetUsers(db *sql.DB, userID int64) (UserList, error) {

	// If no userID is provided (userID is -1), all users are retreived. If
	// a userID is given, a specific user is retreived.
	var rows *sql.Rows

	if userID == -1 {
		queryString := `SELECT * from account`
		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query()
		if err != nil {
			return nil, err
		}

	} else {
		queryString := `SELECT * from account WHERE id = ?`
		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query(userID)
		if err != nil {
			return nil, err
		}
	}

	// Read the rows into the target struct
	var objs UserList

	for rows.Next() {

		var obj User
		err := rows.Scan(
			&obj.ID, &obj.Firstname, &obj.Lastname, &obj.Email,
			&obj.Username, &obj.Password)

		// Forward the error
		if err != nil {
			return nil, err
		}
		// Append object to the list
		objs = append(objs, obj)
	}

	return objs, nil
}
