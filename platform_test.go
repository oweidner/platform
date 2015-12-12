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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oweidner/platform"
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

// TestGetVersion tests GET /platform/version
//
func TestGetVersion(t *testing.T) {

	response, err := http.Get(fmt.Sprintf("%v/platform/version", serverURL))

	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)

	var dat map[string]interface{}
	if err := json.Unmarshal(contents, &dat); err != nil {
		t.Fatal(err)
	}
	v1 := dat["Server"].(string)
	if v1 == "" {
		t.Error("Server version empty.")
	}
	v2 := dat["API"].(string)
	if v2 == "" {
		t.Error("API version empty.")
	}
}

// TestAuth tests POST /platform/auth
