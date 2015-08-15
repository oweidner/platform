//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package picture

import (
	"github.com/codewerft/platform/apiserver/authentication"
	"github.com/codewerft/platform/apiserver/responses"
	"github.com/codewerft/platform/database"

	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
)

// PutPicture updates an account's profile picture.
//
func PutPicture(r render.Render, params martini.Params, db database.Datastore, data ProfilePicture, user authentication.UserInfo) {

	_, err := db.GetDBMap().Exec(`
    UPDATE platform_account SET profile_picture = ? where id = ?`, data.Thumbnail, user.UserID)
	if err != nil {
		responses.Error(r, err.Error())
		return
	}

	responses.OKStatusOnly(r, "Profile picture updated")
}
