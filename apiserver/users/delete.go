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
	"fmt"
	"net/http"
	"strconv"

	"codewerft.net/platform/apiserver/responses"
	"codewerft.net/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// Delete removes one or more user objects from the database and
// sends them back to caller.
//
func Delete(req *http.Request, params martini.Params, r render.Render, db database.Datastore) {

	// userID is either -1 if no user ID was provided or > 0 otherwise.
	var userID int64 = -1

	// Convert the userid string to an integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if params["p1"] != "" {
		var err error
		userID, err = strconv.ParseInt(params["p1"], 10, 64)
		if err != nil {
			responses.DeleteError(r, fmt.Sprintf("Invalid User ID: %v", userID))
			return
		}
	}
	// Use the Generic 'delete' function from the database helper package to
	// delete the user.
	// err := database.GenericDelete(db.Get(), "DELETE FROM user WHERE id = ?", userID)
	// if err != nil {
	// 	responses.DeleteError(r, err.Error())
	// 	return
	// }

	// Respond with OK.
	responses.DeleteOK(r, "User deleted")
}
