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
	"database/sql"
	"fmt"
	"log"
)

// DBWriteLoginError is a convenience function that writes a new login error
// entry to the database.
func DBWriteLoginError(db *sql.DB, origin string, accountname string, details string) {
	// Create the log entry object
	entry := CreateAccessLogEntryRequest{
		Origin:   origin,
		Level:    "ERROR",
		Event:    fmt.Sprintf("Authentication failed: %v", details),
		Username: accountname}
	// Write the entry to the database
	_, err := DBCreateAccessLogEntry(db, entry)
	if err != nil {
		// TODO: handle SQL error on this level.
		log.Printf("%v", err)
	}
}

// DBWriteLoginOK is a convenience function that writes a new login info
// entry to the database.
func DBWriteLoginOK(db *sql.DB, origin string, accountname string) {
	// Create the log entry object
	entry := CreateAccessLogEntryRequest{
		Origin:   origin,
		Level:    "INFO",
		Event:    "Authentication successful.",
		Username: accountname}
	// Write the entry to the database
	_, err := DBCreateAccessLogEntry(db, entry)
	if err != nil {
		// TODO: handle SQL error on this level.
		log.Printf("%v", err)
	}
}

// DBCreateAccessLogEntry creates a new Account object in the database.
//
func DBCreateAccessLogEntry(db *sql.DB, data CreateAccessLogEntryRequest) (AccessLogList, error) {

	stmt, err := db.Prepare(`
		INSERT platform_access_log SET timestamp=NOW(), origin=?, level=?, event=?, accountname=?`)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(
		data.Origin, data.Level, data.Event, data.Username)
	if err != nil {
		return nil, err
	}

	// The id of the newly generated log entry
	logID, _ := res.LastInsertId()
	// Retrieve the newly created object from the database and return it
	logs, err := DBGetLogs(db, logID, -1)
	if err != nil {
		return nil, err
	}
	// Return the account object.
	return logs, nil
}

// DBGetLogs returns a AccessLog object from the database.
func DBGetLogs(db *sql.DB, logID int64, accountID int64) (AccessLogList, error) {

	// If no logID is provided (accountID is -1), all account are retreived. If
	// a logID is given, a specific log entry is retreived.
	var rows *sql.Rows

	queryString := `SELECT * FROM platform_access_log WHERE id = ?`
	// if logID != -1 {
	// 	// If a logID is specified, we only fetch a specific log entry.
	// 	queryString += fmt.Sprintf(`WHERE id = ?`, logID)
	// }
	// 	if accountID != -1 {
	// 		queryString += " AND "
	// 	}
	// }
	// if accountID != -1 {
	// 	// if accountID is specified, we fetch only account-specific logs
	// 	if logID == -1 {
	// 		queryString += " WHERE "
	// 	}
	// 	queryString += fmt.Sprintf("account_id = %v", accountID)
	// }

	stmt, err := db.Prepare(queryString)
	if err != nil {
		return nil, err
	}
	rows, err = stmt.Query(logID)
	if err != nil {
		return nil, err
	}

	// Read the rows into the target struct
	var objs AccessLogList

	for rows.Next() {

		var obj AccessLog
		err := rows.Scan(
			&obj.ID, &obj.Timestamp, &obj.Origin, &obj.Level, &obj.Event,
			&obj.Username)

		// Forward the error
		if err != nil {
			return nil, err
		}
		// Append object to the list
		objs = append(objs, obj)
	}

	return objs, nil
}
