//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package platform_test

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codewerft/platform"
)

var (
	testserver *httptest.Server
	serverURL  string
)

func init() {

	// Config file is passed via -config= flag.
	var configFile = flag.String("config", "", "set configuration file")
	flag.Parse()

	p := platform.New(configFile)

	// ADD APP specific routes and functions

	testserver, serveError := p.UnitTestServe()
	if serveError != nil {
		// report the error
	}

	serverURL = testserver.URL
}

func TestCreateUser(t *testing.T) {

	_, err := http.Get(fmt.Sprintf("%v/platform/version", serverURL))

	if err != nil {
		t.Error(err) //Something is wrong while sending request
	}

	// if request.StatusCode != 200 {
	// 	t.Errorf("Success expected: %d", request.StatusCode) //Uh-oh this means our test failed
	// }
}
