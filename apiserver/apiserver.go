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

	"github.com/codewerft/platform/apiserver/account"
	"github.com/codewerft/platform/apiserver/orgs"
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
	getFn interface{}, listFn interface{}, deleteFn interface{}, createFn interface{}, createReq interface{}, modifyFn interface{}, modifyReq interface{}) {

	/* Add the default list resource
	 */
	r.Get(fmt.Sprintf("%v", basePath),
		strict.Accept("application/json", "text/html"),
		JWTAuth(authEnabled),
		listFn)

	/* Add the default get resource
	 */
	r.Get(fmt.Sprintf("%v/:p1", basePath),
		strict.Accept("application/json", "text/html"),
		JWTAuth(authEnabled),
		getFn)

	/* Add the default delete resource
	 */
	r.Delete(fmt.Sprintf("%v/:p1", basePath),
		strict.Accept("application/json", "text/html"),
		JWTAuth(authEnabled),
		deleteFn)

	/* Add the default create resource
	 */
	r.Put(fmt.Sprintf("%v", basePath),
		strict.Accept("application/json", "text/html"),
		binding.Bind(createReq),
		JWTAuth(authEnabled),
		createFn)

	/* Add the default modify resource
	 */
	r.Post(fmt.Sprintf("%v/:p1", basePath),
		strict.Accept("application/json", "text/html"),
		binding.Bind(modifyReq),
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

	AddDefaultResource(r, fmt.Sprintf("/%v/accounts", prefixPath), authEnabled,
		account.Get, account.Get, account.Delete, account.Create,
		account.CreateAccountRequest{}, account.Modify, account.ModifyAccountRequest{})

	r.Get(fmt.Sprintf("/%v/accesslogs", prefixPath),
		strict.Accept("application/json", "text/html"),
		JWTAuth(authEnabled),
		accesslogs.Get)

	AddDefaultResource(r, fmt.Sprintf("/%v/orgs", prefixPath), authEnabled,
		orgs.Get, orgs.Get, orgs.Delete, orgs.Create,
		orgs.CreateOrgRequest{}, orgs.Modify, orgs.ModifyOrgRequest{})

	AddDefaultResource(r, fmt.Sprintf("/%v/plans", prefixPath), authEnabled,
		plans.Get, plans.Get, plans.Delete, plans.Create,
		plans.CreatePlanRequest{}, plans.Modify, plans.ModifyPlanRequest{})

	AddDefaultResource(r, fmt.Sprintf("/%v/roles", prefixPath), authEnabled,
		roles.Get, roles.Get, roles.Delete, roles.Create,
		roles.CreateRoleRequest{}, roles.Modify, roles.ModifyRoleRequest{})

	//r.NotFound(strict.MethodNotAllowed, strict.NotFound)
	m.Action(r.Handle)

	return m
}
