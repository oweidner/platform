//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package database

import (
	"fmt"

	"gopkg.in/guregu/null.v2"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// HandleStringValue returns a database insert / update pair for a
// string value after checking agains all tags.
func HandleStringValue(value *null.String, tags *[]string, name string) (string, error) {
	// Check if a mandatory value is zero.
	if stringInSlice("req", *tags) && value.IsZero() {
		e := fmt.Errorf("Required field %s is not defined (NULL)", name)
		return "", e
	}
	// Check if a read-only value is set.
	if stringInSlice("!mod", *tags) && !value.IsZero() {
		e := fmt.Errorf("Read-only field %s is set to non-NULL value", name)
		return "", e
	}
	if !value.IsZero() {
		return fmt.Sprintf("%s=%s", (*tags)[0], value.String), nil
	}
	return "", nil
}

// HandleIntValue returns a database insert / update pair for an
// integer value after checking agains all tags.
func HandleIntValue(value *null.Int, tags *[]string, name string) (string, error) {
	// Check if a mandatory value is zero.
	if stringInSlice("req", *tags) && value.IsZero() {
		return "", fmt.Errorf("Required field %s is not defined (NULL)", name)
	}
	// Check if a read-only value is set.
	if stringInSlice("!mod", *tags) && !value.IsZero() {
		return "", fmt.Errorf("Read-only field %s is set to non-NULL value", name)
	}
	if !value.IsZero() {
		return fmt.Sprintf("%s=%d", (*tags)[0], value.Int64), nil
	}
	return "", nil
}

// HandleFloatValue returns a database insert / update pair for a
// float value after checking agains all tags.
func HandleFloatValue(value *null.Float, tags *[]string, name string) (string, error) {
	// Check if a mandatory value is zero.
	if stringInSlice("req", *tags) && value.IsZero() {
		return "", fmt.Errorf("Required field %s is not defined (NULL)", name)
	}
	// Check if a read-only value is set.
	if stringInSlice("!mod", *tags) && !value.IsZero() {
		return "", fmt.Errorf("Read-only field %s is set to non-NULL value", name)
	}
	if !value.IsZero() {
		return fmt.Sprintf("%s=%f", (*tags)[0], value.Float64), nil
	}
	return "", nil
}
