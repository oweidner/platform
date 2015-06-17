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

// New creates a bare bones Platform instance.
func New(configFile *string) *Platform {

	fmt.Printf("test")

	var cfg config.Config

	// Read the configuration.
	err := gcfg.ReadFileInto(&cfg, *configFile)
	if err != nil {
		logging.Log.Fatalf("Error reading configuration file: %v", err)
	}

	// Check configuration semantics
	err = config.CheckConfig(&cfg)
	if err != nil {
		logging.Log.Fatal(err)
	}

	userList := make(map[string]auth.User)

	pwdHash1, _ := bcrypt.GenerateFromPassword([]byte(cfg.Server.AdminPassword), 10)
	userList[cfg.Server.AdminUser] = auth.User{
		Username: cfg.Server.AdminUser,
		Password: pwdHash1,
	}

	jwtPrivateKey, err1 := ioutil.ReadFile(cfg.JWT.PrivateKey)
	if err1 != nil {
		logging.Log.Fatal(fmt.Sprintf("Error reading private key: %v", err1))
	}
	jwtPublicKey, err2 := ioutil.ReadFile(cfg.JWT.PublicKey)
	if err2 != nil {
		logging.Log.Fatal(fmt.Sprintf("Error reading public key: %v", err2))
	}

	// Instantiate / start-up the storage backend
	var ds database.Datastore
	ds = database.NewDefaultDatastore(cfg.MySQL.Host, cfg.MySQL.Database)
	defer ds.Close()

	// Instantiate / start-up the authentication backend
	var ap auth.Authenticator
	ap = auth.NewDefaultAuthProvider(userList)
	defer ap.Close()

	server := apiserver.NewServer(ds, ap, !cfg.Server.DisableAuth, jwtPrivateKey, jwtPublicKey, cfg.JWT.Expiration)

	p := &Platform{Config: &cfg, Server: server}
	return p
}

// Serve launches the Platform HTTP(S) server.
func (p *Platform) Serve() error {

	logging.Log.Info("HTTP server listening on %v", p.Config.Server.Listen)
	if err := http.ListenAndServe(p.Config.Server.Listen, p.Server); err != nil {
		logging.Log.Fatal(fmt.Sprintf("Error starting server: %v", err))
	}

	return nil
}
