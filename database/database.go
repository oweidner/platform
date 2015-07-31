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

	"gopkg.in/gorp.v1"

	"github.com/codewerft/platform/logging"

	_ "github.com/go-sql-driver/mysql"
)

// Datastore provides the common interface for all storage providers.
type Datastore interface {
	GetDB() *sql.DB
	GetDBMap() *gorp.DbMap
	Close()
}

// The DefaultDatastore provides a sqlite-based storage backend.
//
type DefaultDatastore struct {
	DB       sql.DB
	DBMap    gorp.DbMap
	IsClosed bool
}

// NewDefaultDatastore creates a new SQLiteDatastore object.
//
func NewDefaultDatastore(hostname string, dbname string, username string, password string) *DefaultDatastore {

	var connect = fmt.Sprintf("%v:%v@%v/%v?charset=utf8&parseTime=true", username, password, hostname, dbname)
	// connectCleaned ommits the password from the string for security reasons
	var connectCleaned = fmt.Sprintf("%v:%v@%v/%v?charset=utf8&parseTime=true", username, "***", hostname, dbname)

	logging.Log.Info("Connecting to MySQL server as %v", connectCleaned)

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

	// Set up GORP
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	ds := DefaultDatastore{
		// Session:  *session,
		DB:       *db,
		DBMap:    *dbmap,
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

// GetDB returns the database handle.
//
func (ds *DefaultDatastore) GetDB() *sql.DB {
	// ds.Session.Close()
	return &ds.DB
}

// GetDBMap returns the database map handle.
//
func (ds *DefaultDatastore) GetDBMap() *gorp.DbMap {
	// ds.Session.Close()
	return &ds.DBMap
}
