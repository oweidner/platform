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
	"database/sql"
	"fmt"

	"github.com/codewerft/platform/logging"

	_ "github.com/go-sql-driver/mysql"
)

// Datastore provides the common interface for all storage providers.
type Datastore interface {
	Get() *sql.DB
	Close()
}

// The DefaultDatastore provides a sqlite-based storage backend.
//
type DefaultDatastore struct {
	DB       sql.DB
	IsClosed bool
}

// NewDefaultDatastore creates a new SQLiteDatastore object.
//
func NewDefaultDatastore(hostname string, dbname string, accountname string, password string) *DefaultDatastore {

	var connect = fmt.Sprintf("%v:%v@%v/%v?charset=utf8&parseTime=true", accountname, password, hostname, dbname)
	// connectCleaned ommits the password from the string for security reasons
	var connectCleaned = fmt.Sprintf("%v:%v@%v/%v?charset=utf8&parseTime=true", accountname, "***", hostname, dbname)

	logging.Log.Info("Connecting to MySQL server %v", connectCleaned)

	db, err := sql.Open("mysql", connect)
	if err != nil {
		logging.Log.Fatalf("Connection to MySQL server failed: %v", err)
	}

	// Ping() asserts that the conenction was established successfully.
	err = db.Ping()
	if err != nil {
		logging.Log.Fatalf("Connection to MySQL database failed: %v", err)
	}
	logging.Log.Info("Connection to MySQL server established.")

	ds := DefaultDatastore{
		// Session:  *session,
		DB:       *db,
		IsClosed: false,
	}
	return &ds
}

// Close implements the Datastore interface.
//
func (ds *DefaultDatastore) Close() {
	// ds.Session.Close()
	ds.IsClosed = true
}

// Close implements the Datastore interface.
//
func (ds *DefaultDatastore) Get() *sql.DB {
	// ds.Session.Close()
	return &ds.DB
}
