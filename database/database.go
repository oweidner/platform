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
func NewDefaultDatastore(hostname string, dbname string) *DefaultDatastore {

	db, err := sql.Open("mysql", "root:@/ohoi.dev?charset=utf8")
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
