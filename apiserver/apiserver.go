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

	"github.com/codewerft/platform/apiserver/orgs"
	"github.com/codewerft/platform/apiserver/plans"
	"github.com/codewerft/platform/apiserver/users"

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
	r.Put(fmt.Sprintf("%v/:p1", basePath),
		strict.Accept("application/json", "text/html"),
		binding.Bind(modifyReq),
		JWTAuth(authEnabled),
		createFn)
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

	AddDefaultResource(r, fmt.Sprintf("/%v/users", prefixPath), authEnabled,
		users.Get, users.Get, users.Delete, users.Create,
		users.CreateUserRequest{}, users.Modify, users.ModifyUserRequest{})

	AddDefaultResource(r, fmt.Sprintf("/%v/orgs", prefixPath), authEnabled,
		orgs.Get, orgs.Get, orgs.Delete, orgs.Create,
		orgs.CreateOrgRequest{}, orgs.Modify, orgs.ModifyOrgRequest{})

	AddDefaultResource(r, fmt.Sprintf("/%v/plans", prefixPath), authEnabled,
		plans.Get, plans.Get, plans.Delete, plans.Create,
		plans.CreatePlanRequest{}, plans.Modify, plans.ModifyPlanRequest{})

	// // -------------------- The /plans resource --------------------
	//
	// /* Retreive a list of all available plans.
	//  */
	// r.Get("/plans",
	// 	strict.Accept("application/json", "text/html"),
	// 	JWTAuth(authEnabled),
	// 	plans.GetPlans)
	//
	// /* Retreive information about a specific plan.
	//  */
	// r.Get("/plans/:pid",
	// 	strict.Accept("application/json", "text/html"),
	// 	JWTAuth(authEnabled),
	// 	users.GetPlans)
	//
	// /* Create a new plan.
	//  */
	// r.Put("/plans/:pid",
	// 	strict.Accept("application/json", "text/html"),
	// 	binding.Bind(users.CreateUserRequest{}),
	// 	JWTAuth(authEnabled),
	// 	users.CreatePlan)
	//
	// /* Delete a plan.
	//  */
	// r.Delete("/plans/:pid",
	// 	strict.Accept("application/json", "text/html"),
	// 	JWTAuth(authEnabled),
	// 	users.DeletePlan)

	//r.NotFound(strict.MethodNotAllowed, strict.NotFound)
	m.Action(r.Handle)

	return m
}
