//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package organizations

import "database/sql"

// DBGetOrganizations returns a list of User object from the database.
//
func DBGetOrganizations(db *sql.DB, orgID int64) ([]Organization, error) {

	// If no userID is provided (userID is -1), all users are retreived. If
	// a userID is given, a specific user is retreived.
	var rows *sql.Rows

	if orgID == -1 {
		queryString := `SELECT * from organization`
		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query()
		if err != nil {
			return nil, err
		}

	} else {
		queryString := `SELECT * from organization WHERE id = ?`
		stmt, err := db.Prepare(queryString)
		if err != nil {
			return nil, err
		}
		rows, err = stmt.Query(orgID)
		if err != nil {
			return nil, err
		}
	}

	// Read the rows into the target struct
	var objs []Organization
	//
	rows, err := stmt.Query()
	for rows.Next() {

		var obj Organization
		err = rows.Scan(
			&obj.ID, &obj.Name, &obj.LocationLat, &obj.LocationLon,
			&obj.ContactStreet, &obj.ContactCity, &obj.ContactZip,
			&obj.ContactCountry, &obj.ContactPhone, &obj.ContactEmail)

		// Forward the error
		if err != nil {
			return []Organization{}, err
		}
		// Append object to the list
		objs = append(objs, obj)
	}

	return objs, nil
}
