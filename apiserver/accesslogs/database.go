//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package accesslogs

import (
	"fmt"
	"log"
	"time"

	"github.com/codewerft/platform/database"

	"gopkg.in/guregu/null.v2"
)

// DBWriteLoginError is a convenience function that writes a new login error
// entry to the database.
func DBWriteLoginError(db database.Datastore, origin string, username string, details string, userid null.Int) {

	// Create the log entry object
	entry := AccessLogEntry{
		Origin:    null.StringFrom(origin),
		Timestamp: time.Now().UTC(),
		Level:     null.StringFrom("ERROR"),
		Event:     null.StringFrom(fmt.Sprintf("Authentication failed for user %v: %v", username, details)),
		Username:  null.StringFrom(username),
		UserID:    userid}

	// Insert into the database
	err := db.GetDBMap().Insert(&entry)
	if err != nil {
		log.Printf("%v", err)
	}
}

// DBWriteLoginOK is a convenience function that writes a new login info
// entry to the database.
func DBWriteLoginOK(db database.Datastore, origin string, username string, userid null.Int) {
	// Create the log entry object
	entry := AccessLogEntry{
		Origin:    null.StringFrom(origin),
		Timestamp: time.Now().UTC(),
		Level:     null.StringFrom("INFO"),
		Event:     null.StringFrom(fmt.Sprintf("Authentication granted for user %v", username)),
		Username:  null.StringFrom(username),
		UserID:    userid}

	// Insert into the database
	err := db.GetDBMap().Insert(&entry)
	if err != nil {
		log.Printf("%v", err)
	}
}
