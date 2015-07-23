//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package apiserver

import (
	"fmt"

	"github.com/codewerft/platform/auth"
	"github.com/codewerft/platform/database"
	"github.com/codewerft/platform/logging"

	"github.com/codewerft/platform/apiserver/accounts"
	"github.com/codewerft/platform/apiserver/organisations"
	"github.com/codewerft/platform/apiserver/plans"
	"github.com/codewerft/platform/apiserver/roles"

	"github.com/codewerft/platform/apiserver/accesslogs"

	"github.com/attilaolah/strict"
	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/cors"
)

// The resource path prefix. If not blank, it requires a '/' as prefix.
const APIPrefix = ""

var (
	jwtExpiration int
	jwtPrivateKey []byte
	jwtPublicKey  []byte
)

type getFunc func(int) string

func AddDefaultResource(r martini.Router, basePath string, authEnabled bool,
	listFn interface{}, getFn interface{}, deleteFn interface{}, createFn interface{}, modifyFn interface{}, datatype interface{}) {

	/* List the resources
	 */
	r.Get(fmt.Sprintf("%v/:p1", basePath),
		strict.Accept("application/json", "text/html"),
		JWTAuth(authEnabled),
		getFn)

	/* Get a resource
	 */
	r.Get(fmt.Sprintf("%v", basePath),
		strict.Accept("application/json", "text/html"),
		JWTAuth(authEnabled),
		listFn)

	/* Delete a resource
	 */
	r.Delete(fmt.Sprintf("%v/:p1", basePath),
		strict.Accept("application/json", "text/html"),
		JWTAuth(authEnabled),
		deleteFn)

	/* Create a new resource
	 */
	r.Put(fmt.Sprintf("%v", basePath),
		strict.Accept("application/json", "text/html"),
		binding.Bind(datatype),
		JWTAuth(authEnabled),
		createFn)

	/* Modify a resource
	 */
	r.Post(fmt.Sprintf("%v/:p1", basePath),
		strict.Accept("application/json", "text/html"),
		binding.Bind(datatype),
		JWTAuth(authEnabled),
		modifyFn)
}

// NewServer creates a new Server instance.
func NewServer(ds database.Datastore, ap auth.Authenticator, prefixPath string, authEnabled bool,
	privKey []byte, pubKey []byte, expiration int) *martini.Martini {

	// Print a big fat warning if authentication is disabled.
	if !authEnabled {
		logging.Log.Warning("***********************************************************************************")
		logging.Log.Warning("!! AUTHENTICATION DISABLED -- DO NOT RUN THIS SERVER IN A PRODUCTION ENVIRONMENT !!")
		logging.Log.Warning("***********************************************************************************")
	}

	jwtPrivateKey = privKey
	jwtPublicKey = pubKey

	// Export JWT expiration
	jwtExpiration = expiration
	logging.Log.Info("Server config: JWT expiration set to %v hours", expiration)

	// Setup middleware
	var m *martini.Martini

	m = martini.New()
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	m.Use(render.Renderer())

	// Configure CORS (Cross-origin resource sharing). Without this, calling
	// this REST API from within an Angular app (e.g. via $http.get) will
	// most definetly fail. See Wikipedia for details:
	// http://en.wikipedia.org/wiki/Cross-origin_resource_sharing
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Inject authenticator
	m.MapTo(ap, (*auth.Authenticator)(nil))

	// Inject datastore
	m.MapTo(ds, (*database.Datastore)(nil))
	// Add the router action

	r := martini.NewRouter()

	r.Get(fmt.Sprintf("/%v/version", prefixPath),
		strict.Accept("application/json", "text/html"),
		JWTAuth(authEnabled && false),
		GetVersion)

	// Auth API
	r.Post("/auth",
		strict.Accept("application/json", "text/html"),
		binding.Bind(AuthRequest{}),
		Auth)

	// User info API
	r.Get("/accounts/me",
		strict.Accept("application/json", "text/html"),
		JWTAuth(authEnabled && false),
		GetSelf)

	r.Get(fmt.Sprintf("/%v/accesslogs", prefixPath),
		strict.Accept("application/json", "text/html"),
		JWTAuth(authEnabled),
		accesslogs.Get)

	// Defines /resources/* resources
	//
	AddDefaultResource(r, fmt.Sprintf("/%v/accounts", prefixPath), authEnabled,
		accounts.List, accounts.Get, accounts.Delete, accounts.Create, accounts.Modify,
		accounts.Account{})

	// Defines /organisations/* resources
	//
	AddDefaultResource(r, fmt.Sprintf("/%v/organisations", prefixPath), authEnabled,
		organisations.List, organisations.Get, organisations.Delete, organisations.Create, organisations.Modify,
		organisations.Organisation{})

	// Defines /plans/* resources
	//
	AddDefaultResource(r, fmt.Sprintf("/%v/plans", prefixPath), authEnabled,
		plans.List, plans.Get, plans.Delete, plans.Create, plans.Modify,
		plans.Plan{})

	// Defines /roles/* resources
	//
	AddDefaultResource(r, fmt.Sprintf("/%v/roles", prefixPath), authEnabled,
		roles.List, roles.Get, roles.Delete, roles.Create, roles.Modify,
		roles.Role{})

	//r.NotFound(strict.MethodNotAllowed, strict.NotFound)
	m.Action(r.Handle)

	return m
}
