//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package main

import (
	"flag"
	"net/http"

	"github.com/oweidner/platform"
	"github.com/oweidner/platform/database"
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
)

type Test struct {
	Name string
}

// GetVersion returns the version of Codewerft Platform.
func Tryme(r render.Render, params martini.Params, db database.Datastore) {
	r.JSON(http.StatusOK, "{'X': 'Y'}")
	return
}

// p is the global Platform instance.
var p *platform.Platform

func main() {

	// Config file is passed via -config= flag.
	var configFile = flag.String("config", "", "set configuration file")
	flag.Parse()

	p = platform.New(configFile)

	p.AddGORPTable("tutorbox_students", "id", Test{})
	p.Get("/test", Tryme)

	// ADD APP specific routes and functions

	serveError := p.Serve()
	if serveError != nil {
		// report the error
	}
}
