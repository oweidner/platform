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
)

// DBCreateAccessLogEntry creates a new User object in the database.
//
func DBCreateAccessLogEntry(db *sql.DB, data CreateAccessLogEntryRequest) (AccessLogList, error) {

	fmt.Printf("SDASDASDS")

	stmt, err := db.Prepare(`
		INSERT access_log SET timestamp=NOW(), origin=?, level=?, event=?, account_id=?`)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(
		data.Origin, data.Level, data.Event, data.AccountID)
	if err != nil {
		return nil, err
	}

	// The id of the newly generated log entry
	logID, _ := res.LastInsertId()
	// Retrieve the newly created object from the database and return it
	logs, err := DBGetLogs(db, logID)
	if err != nil {
		return nil, err
	}
	// Return the user object.
	return logs, nil
}

// DBGetLogs returns a AccessLog object from the database.
func DBGetLogs(db *sql.DB, logID int64) (AccessLogList, error) {

	// If no logID is provided (userID is -1), all users are retreived. If
	// a logID is given, a specific log entry is retreived.
	var rows *sql.Rows

	if logID == -1 {
		queryString := `SELECT * from access_log`
		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query()
		if err != nil {
			return nil, err
		}

	} else {
		queryString := `SELECT * from access_log WHERE id = ?`
		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query(logID)
		if err != nil {
			return nil, err
		}
	}

	// Read the rows into the target struct
	var objs AccessLogList

	for rows.Next() {

		var obj AccessLog
		err := rows.Scan(
			&obj.ID, &obj.Timestamp, &obj.Origin, &obj.Level, &obj.Event,
			&obj.AccountID)

		// Forward the error
		if err != nil {
			return nil, err
		}
		// Append object to the list
		objs = append(objs, obj)
	}

	return objs, nil
}
