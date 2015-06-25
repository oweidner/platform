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

	"code.google.com/p/gcfg"
	"github.com/codewerft/platform/apiserver"
	"github.com/codewerft/platform/apiserver/users"

	"github.com/codewerft/platform/auth"
	"github.com/codewerft/platform/config"
	"github.com/codewerft/platform/database"
	"github.com/codewerft/platform/logging"
	"github.com/go-martini/martini"
	"golang.org/x/crypto/bcrypt"
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

	var cfg config.Config

	// Read the configuration.
	err := gcfg.ReadFileInto(&cfg, *configFile)
	if err != nil {
		logging.Log.Fatalf("Error reading configuration file: %v", err)
	}

	// Check configuration semantics
	err = config.CheckConfig(&cfg, *configFile)
	if err != nil {
		logging.Log.Fatal(err)
	}

	// Create the root user credentials from the username and password
	// values defined in the config file.
	rootUser := users.User{}
	pwdHash1, _ := bcrypt.GenerateFromPassword([]byte(cfg.Server.AdminPassword), 0)
	rootUser = users.User{
		Firstname: "Root",
		Lastname:  "Admin User",
		Username:  cfg.Server.AdminUser,
		Password:  string(pwdHash1)}
	// Load the JWT __PRIVATE__ key from the path / filename defined in
	// the config file.
	jwtPrivateKey, err1 := ioutil.ReadFile(cfg.JWT.PrivateKey)
	if err1 != nil {
		logging.Log.Fatal(fmt.Sprintf("Error reading private key: %v", err1))
	}
	// Load the JWT __PUBLIC__ key from the path / filename defined in
	// the config file.
	jwtPublicKey, err2 := ioutil.ReadFile(cfg.JWT.PublicKey)
	if err2 != nil {
		logging.Log.Fatal(fmt.Sprintf("Error reading public key: %v", err2))
	}
	// Instantiate the storage database backend with the values defined
	// in the config file.
	ds = database.NewDefaultDatastore(cfg.MySQL.Host, cfg.MySQL.Database, cfg.MySQL.Username, cfg.MySQL.Password)
	defer ds.Close()
	// Instantiate the authentication backend and inject the root user.
	ap = auth.NewDefaultAuthProvider(ds, rootUser)

	// Finally, we start up the Platform API server and inject the storage
	// and authentication backend instances.
	server := apiserver.NewServer(ds, ap, cfg.Server.APIPrefix,
		!cfg.Server.DisableAuth,
		jwtPrivateKey, jwtPublicKey, cfg.JWT.Expiration)

	// Return the Platform handle.
	return &Platform{Config: &cfg, Server: server}
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
