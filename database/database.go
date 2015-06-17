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
func NewDefaultDatastore(hostname string, dbname string, username string, password string) *DefaultDatastore {

	var connect = fmt.Sprintf("%v:%v@%v/%v?charset=utf8", username, password, hostname, dbname)
	// connectCleaned ommits the password from the string for security reasons
	var connectCleaned = fmt.Sprintf("%v:%v@%v/%v?charset=utf8", username, "***", hostname, dbname)

	logging.Log.Info("Connecting to MySQL: %v", connectCleaned)

	db, err := sql.Open("mysql", connect)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)

	}

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
