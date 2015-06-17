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
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/codewerft/platform/logging"
)

func logError(err error) error {

	_, file, lineno, _ := runtime.Caller(1)
	loc := fmt.Sprintf("%v:%v", filepath.Base(file), lineno)

	if err != nil {
		logging.Log.Error(fmt.Sprintf("%v (%v)", err, loc))
		return err
	}
	return nil
}

func handleError(err error, msg string) error {

	_, file, lineno, _ := runtime.Caller(1)
	loc := fmt.Sprintf("%v:%v", filepath.Base(file), lineno)

	if err != nil {
		se := NewDataStoreError(msg, err.Error(), "Storage Layer", loc)
		logging.Log.Error(fmt.Sprintf("%v", se))
		return se
	}
	return nil
}

func handleFatalErrorAndExit(err error, msg string) {
	_, file, lineno, _ := runtime.Caller(1)
	loc := fmt.Sprintf("%v:%v", filepath.Base(file), lineno)

	if err != nil {
		se := NewDataStoreError(msg, err.Error(), "Storage Layer", loc)
		log.Fatalf("%v", se)
	}
}

// DataStoreError is an error type emitted by the storage layer.
type DataStoreError struct {
	// when did the error occur?
	When time.Time
	// where did the error occur?
	Where string
	// details of where did the error occur?
	WhereDetails string
	// what is the error?
	What string
	// what are the details of the error?
	WhatDetails string
}

// NewDataStoreError creates a new NewDataStoreError object.
func NewDataStoreError(what string, whatDetails string, where string, whereDetails string) DataStoreError {

	return DataStoreError{
		When:         time.Now().Local(),
		Where:        where,
		WhereDetails: whereDetails,
		What:         what,
		WhatDetails:  whatDetails}
}

// String representation of the error
func (e DataStoreError) Error() string {
	return fmt.Sprintf("[%v] (%v): %v (%v)", e.Where, e.WhereDetails, e.What, e.WhatDetails)
}
