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
	"image"
	"image/png"
	"io"

	"github.com/oweidner/platform/apiserver/authentication"
	"github.com/oweidner/platform/apiserver/responses"
	"github.com/oweidner/platform/database"

	"github.com/martini-contrib/render"

	// Available image format drivers
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// convertToPNG converts from any recognized format to PNG.
func convertToPNG(w io.Writer, r io.Reader) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	return png.Encode(w, img)
}

// GetPicture returns the profile picture data for an account.
//
func GetPicture(r render.Render, db database.Datastore, user authentication.UserInfo) {

	var picture ProfilePicture

	dbError := db.GetDBMap().SelectOne(&picture, `
    SELECT profile_picture from platform_account WHERE id = ?;`, user.UserID)
	if dbError != nil {
		responses.Error(r, dbError.Error())
		return
	}

	responses.OKStatusPlusData(r, picture, 1)
}
