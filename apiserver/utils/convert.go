//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package utils

import (
	"strconv"

	"codewerft.net/ohoi/apiserver/admin/apiserver/responses"
	"github.com/gavv/martini-render"
)

// StringToInt64 converts a string to 64-bit integer nubmer. If the
// conversion was successful, the int64 number is returned. If not,
// an error response is sent back to the caller.
func StringToInt64(r render.Render, number string) int64 {
	accountID, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		responses.GetError(r, "Invalid AccountID.")
	}
	return accountID
}
