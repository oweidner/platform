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
	"fmt"
	"strconv"
)

// ParseResourceID converts a string to 64-bit integer nubmer. If the
// conversion was successful, the int64 number is returned. If not,
// an error response is sent back to the caller.
func ParseResourceID(idString string) (int64, error) {
	// planID is either -1 if no plan ID was provided or > 0 otherwise.
	var theID int64 = -1

	// Convert the ID string to a 64-bit integer. In case the conversion
	// fails, an error response is sent back to the caller.
	if idString != "" {
		var err error
		theID, err = strconv.ParseInt(idString, 10, 64)
		if err != nil {
			return theID, fmt.Errorf("Invalid resource ID: %v", idString)
		}
	}
	return theID, nil
}
