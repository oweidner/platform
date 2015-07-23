//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package platform

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"

	"code.google.com/p/gcfg"
	"github.com/codewerft/platform/apiserver"
	"github.com/codewerft/platform/apiserver/accounts"

	"github.com/codewerft/platform/auth"
	"github.com/codewerft/platform/config"
	"github.com/codewerft/platform/database"
	"github.com/codewerft/platform/logging"
	"github.com/go-martini/martini"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ServerVersion holds the version of the platform server.
	// Set externally via -ldflags "-X platform.ServerVersion x.y.z"
	ServerVersion string

	// APIVersion holds the version of the API.
	// Set externally via -ldflags "-X platform.APIVersion x.y.z"
	APIVersion string
)

// Platform represents the top level application.
type Platform struct {
	Router *martini.Router
	Config *config.Config
	Server *martini.Martini
}

var ap auth.Authenticator
var ds database.Datastore

// New creates a bare bones Platform instance.
func New(configFile *string) *Platform {

	// get the base path of the configuration file
	configFileAbs, absErr := filepath.Abs(*configFile)
	if absErr != nil {
		// file doesn't exist
		logging.Log.Fatal(fmt.Sprintf("Couldn't open configuration file %v", absErr))
	}
	basePath := path.Dir(configFileAbs)

	var cfg config.Config

	// Read the configuration.
	err := gcfg.ReadFileInto(&cfg, configFileAbs)
	if err != nil {
		logging.Log.Fatalf("Error reading configuration file: %v", err)
	}

	// Check configuration semantics
	err = config.CheckConfig(&cfg, configFileAbs, basePath)
	if err != nil {
		logging.Log.Fatal(err)
	}

	// Create the root account credentials from the username and password
	// values defined in the config file.
	rootAccount := accounts.Account{}
	pwdHash1, _ := bcrypt.GenerateFromPassword([]byte(cfg.Server.AdminPassword), 0)
	rootAccount = accounts.Account{
		Firstname: "Root",
		Lastname:  "Admin Account",
		Username:  cfg.Server.AdminAccount,
		Password:  string(pwdHash1)}
	// Load the JWT __PRIVATE__ key from the path / filename defined in
	// the config file.
	jwtPrivateKey, err1 := ioutil.ReadFile(cfg.JWT.PrivateKey)
	if err1 != nil {
		logging.Log.Fatal(fmt.Sprintf("Error reading private key: %v", err1))
	}
	logging.Log.Info(fmt.Sprintf("Loaded JWT private key from %v", cfg.JWT.PrivateKey))

	// Load the JWT __PUBLIC__ key from the path / filename defined in
	// the config file.
	jwtPublicKey, err2 := ioutil.ReadFile(cfg.JWT.PublicKey)
	if err2 != nil {
		logging.Log.Fatal(fmt.Sprintf("Error reading public key: %v", err2))
	}
	logging.Log.Info(fmt.Sprintf("Loaded JWT public key from %v", cfg.JWT.PublicKey))

	// Instantiate the storage database backend with the values defined
	// in the config file.
	ds = database.NewDefaultDatastore(cfg.MySQL.Host, cfg.MySQL.Database, cfg.MySQL.Username, cfg.MySQL.Password)
	defer ds.Close()
	// Instantiate the authentication backend and inject the root account.
	ap = auth.NewDefaultAuthProvider(ds, rootAccount)

	// Finally, we start up the Platform API server and inject the storage
	// and authentication backend instances.
	server := apiserver.NewServer(ds, ap, cfg.Server.APIPrefix,
		!cfg.Server.DisableAuth,
		jwtPrivateKey, jwtPublicKey, cfg.JWT.Expiration)

	// Return the Platform handle.
	return &Platform{Config: &cfg, Server: server}
}

// UnitTestServe launches the Platform HTTPS server for running the unit tests.
func (p *Platform) UnitTestServe() (*httptest.Server, error) {
	server := httptest.NewServer(p.Server)
	logging.Log.Info("Codewerft Platform unit test server running on %s", server.URL)

	return server, nil
}

// Serve launches the Platform HTTP(S) server.
func (p *Platform) Serve() error {
	// if TLS is enable in the configuration file, we start
	// an HTTPS server with the provided X.509 certificates,
	// otherwise, start an HTTP server.without TLS.
	if p.Config.TLS.EnableTLS == true {
		logging.Log.Info("HTTPS/TLS enabled. Using X.509 keypair %v and %v", p.Config.TLS.CertFile, p.Config.TLS.KeyFile)
		logging.Log.Info("Codewerft Platform server available at https://localhost%v", p.Config.Server.Listen)

		if err := http.ListenAndServeTLS(
			p.Config.Server.Listen,
			p.Config.TLS.CertFile,
			p.Config.TLS.KeyFile,
			p.Server); err != nil {
			logging.Log.Fatalf("Error starting Codewerft Platform server: %v", err)
		}

	} else {
		logging.Log.Warning("***********************************************************************************")
		logging.Log.Warning("!! HTTPS/TLS DISABLED -- DO NOT RUN THIS SERVER IN A PRODUCTION ENVIRONMENT      !!")
		logging.Log.Warning("***********************************************************************************")

		logging.Log.Info("Codewerft Platform server available at http://localhost%v", p.Config.Server.Listen)
		if err := http.ListenAndServe(p.Config.Server.Listen, p.Server); err != nil {
			logging.Log.Fatal(fmt.Sprintf("Error starting Codewerft Platform server: %v", err))
		}
	}
	return nil
}
