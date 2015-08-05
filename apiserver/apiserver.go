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

	"github.com/codewerft/platform/database"
	"github.com/codewerft/platform/logging"

	"github.com/codewerft/platform/apiserver/accounts"
	"github.com/codewerft/platform/apiserver/accounts/accountpassword"
	"github.com/codewerft/platform/apiserver/accounts/accountplans"
	"github.com/codewerft/platform/apiserver/accounts/accountroles"
	"github.com/codewerft/platform/apiserver/accounts/accountstatus"
	"github.com/codewerft/platform/apiserver/authentication"
	"github.com/codewerft/platform/apiserver/organisations"
	"github.com/codewerft/platform/apiserver/organisations/orgplans"
	"github.com/codewerft/platform/apiserver/organisations/orgstatus"

	"github.com/codewerft/platform/apiserver/plans"
	"github.com/codewerft/platform/apiserver/roles"

	"github.com/codewerft/platform/apiserver/accesslogs"

	"github.com/attilaolah/strict"
	"github.com/gavv/martini-render"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/cors"
)

// The resource path prefixes
const PlatformPrefix = ""
const ApplicationPrefix = ""

// The system roles
const PlatformAdminRole = "PLATFORM_ADMIN"
const PlatformUserRole = "PLATFORM_USER"

func AddDefaultResource(r martini.Router, basePath string, idPlaceholder string, jwtcfg authentication.JWTConfig,
	listFn interface{}, getFn interface{}, deleteFn interface{}, createFn interface{}, modifyFn interface{}, datatype interface{}) {

	/* List the resources
	 */
	if getFn != nil {
		r.Get(fmt.Sprintf("%v/:%s", basePath, idPlaceholder),
			strict.Accept("application/json", "text/html"),
			authentication.JWTAuth(jwtcfg, PlatformAdminRole),
			getFn)
	}

	/* Get a resource
	 */
	if listFn != nil {
		r.Get(fmt.Sprintf("%v", basePath),
			strict.Accept("application/json", "text/html"),
			authentication.JWTAuth(jwtcfg, PlatformAdminRole),
			listFn)
	}

	/* Delete a resource
	 */
	if deleteFn != nil {
		r.Delete(fmt.Sprintf("%v/:%s", basePath, idPlaceholder),
			strict.Accept("application/json", "text/html"),
			authentication.JWTAuth(jwtcfg, PlatformAdminRole),
			deleteFn)
	}

	/* Create a new resource
	 */
	if createFn != nil {
		r.Put(fmt.Sprintf("%v", basePath),
			strict.Accept("application/json", "text/html"),
			binding.Bind(datatype),
			authentication.JWTAuth(jwtcfg, PlatformAdminRole),
			createFn)
	}

	/* Modify a resource
	 */
	if modifyFn != nil {
		r.Post(fmt.Sprintf("%v/:%s", basePath, idPlaceholder),
			strict.Accept("application/json", "text/html"),
			binding.Bind(datatype),
			authentication.JWTAuth(jwtcfg, PlatformAdminRole),
			modifyFn)
	}
}

// Platform represents the top level application.
type Server struct {
	Name      string
	Martini   *martini.Martini
	Router    martini.Router
	JWTConfig authentication.JWTConfig
}

// // New creates Server instance
// func New(ds database.Datastore, ap auth.Authenticator, prefixPath string,
// 	authEnabled bool, platformEnabled bool, privKey []byte, pubKey []byte, expiration int) *Server {
// 	return &Server{Name: "Server"}
// }

func (p *Server) Get(path string, fn interface{}) error {
	return nil
}

// // NewServer creates a new Server instance.
// func NewServer(ds database.Datastore, ap auth.Authenticator, prefixPath string,
// 	authEnabled bool, platformEnabled bool, privKey []byte, pubKey []byte, expiration int) *martini.Martini {

// New creates Server instance
func New(ds database.Datastore, prefixPath string,
	authEnabled bool, platformEnabled bool, privKey []byte, pubKey []byte, expiration int) *Server {

	// Print a big fat warning if authentication is disabled.
	if !authEnabled {
		logging.Log.Warning("***********************************************************************************")
		logging.Log.Warning("!! AUTHENTICATION DISABLED -- DO NOT RUN THIS SERVER IN A PRODUCTION ENVIRONMENT !!")
		logging.Log.Warning("***********************************************************************************")
	}

	jwtcfg := authentication.JWTConfig{
		PublicKey:  pubKey,
		PrivateKey: privKey,
		Expiration: expiration}

	logging.Log.Info("Server config: JWT expiration set to %v hours", expiration)

	// Configure GORP
	ds.GetDBMap().AddTableWithName(accesslogs.AccessLogEntry{}, "platform_access_log").SetKeys(true, "id")

	ds.GetDBMap().AddTableWithName(accounts.Account{}, "platform_account").SetKeys(true, "id")
	ds.GetDBMap().AddTableWithName(accountroles.AccountOrganisationRole{}, "platform_account_organisation_role").SetKeys(true, "id")
	ds.GetDBMap().AddTableWithName(accountplans.AccountPlanAssoc{}, "platform_account_plan_assoc").SetKeys(true, "id")

	ds.GetDBMap().AddTableWithName(organisations.Organisation{}, "platform_organisation").SetKeys(true, "id")
	ds.GetDBMap().AddTableWithName(orgplans.OrganisationPlanAssoc{}, "platform_organisation_plan_assoc").SetKeys(true, "id")

	ds.GetDBMap().AddTableWithName(plans.Plan{}, "platform_plan").SetKeys(true, "id")
	ds.GetDBMap().AddTableWithName(roles.Role{}, "platform_role").SetKeys(true, "id")

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

	// Inject datastore
	m.MapTo(ds, (*database.Datastore)(nil))

	var jwtc authentication.JWTConfigIF
	jwtc = authentication.TheJWTConfig{Config: jwtcfg}

	m.MapTo(jwtc, (*authentication.JWTConfigIF)(nil))

	r := martini.NewRouter()

	r.Get(fmt.Sprintf("/%v/version", prefixPath),
		strict.Accept("application/json", "text/html"),
		authentication.JWTAuth(jwtcfg, PlatformAdminRole),
		GetVersion)

	// Auth API
	r.Post("/auth",
		strict.Accept("application/json", "text/html"),
		binding.Bind(authentication.AuthRequest{}),
		authentication.Auth)

	// User info API
	r.Get("/self",
		strict.Accept("application/json", "text/html"),
		authentication.GetSelf)

	// Change own data
	// r.Post("/self",
	// 	strict.Accept("application/json", "text/html"),
	// 	JWTAuth(false, PlatformAdminRole),
	// 	GetSelf)

	// Change own password
	// r.Post("/self/passsword",
	// 	strict.Accept("application/json", "text/html"),
	// 	JWTAuth(false, PlatformAdminRole),
	// 	GetSelf)

	// We add the platform API *ONLY* if it was enabled in the
	// configuration file.
	if platformEnabled == true {

		r.Get(fmt.Sprintf("/%v/accesslogs", prefixPath),
			strict.Accept("application/json", "text/html"),
			authentication.JWTAuth(jwtcfg, PlatformAdminRole),
			accesslogs.Get)

		// Defines /accounts/* resources
		//
		AddDefaultResource(r, fmt.Sprintf("/%v/accounts", prefixPath), "p1", jwtcfg,
			accounts.List, accounts.Get, accounts.Delete, accounts.Create, accounts.Modify,
			accounts.Account{})

		// Defines /accounts/*/password resources
		//
		r.Post(fmt.Sprintf("/%v/accounts/:p1/password", prefixPath),
			strict.Accept("application/json", "text/html"),
			binding.Bind(accountpassword.AccountPassword{}),
			authentication.JWTAuth(jwtcfg, PlatformAdminRole),
			accountpassword.Set)

		// Defines /accounts/*/status resources
		//
		r.Post(fmt.Sprintf("/%v/accounts/:p1/status", prefixPath),
			strict.Accept("application/json", "text/html"),
			binding.Bind(accountstatus.AccountStatus{}),
			authentication.JWTAuth(jwtcfg, PlatformAdminRole),
			accountstatus.Set)

		// Defines /accounts/*/plans resources
		//
		AddDefaultResource(r, fmt.Sprintf("/%v/accounts/:p1/plans", prefixPath), "p2", jwtcfg,
			accountplans.List, accountplans.Get, nil, accountplans.Create, nil,
			accountplans.AccountPlanAssoc{})

		// Defines /accounts/*/roles resources
		//
		AddDefaultResource(r, fmt.Sprintf("/%v/accounts/:p1/roles", prefixPath), "p2", jwtcfg,
			accountroles.List, accountroles.Get, nil, accountroles.Create, nil,
			accountroles.AccountOrganisationRole{})

		// Defines /organisations/* resources
		//
		AddDefaultResource(r, fmt.Sprintf("/%v/organisations", prefixPath), "p1", jwtcfg,
			organisations.List, organisations.Get, organisations.Delete, organisations.Create, organisations.Modify,
			organisations.Organisation{})

		// Defines /organisations/*/status resources
		//
		r.Post(fmt.Sprintf("/%v/organisations/:p1/status", prefixPath),
			strict.Accept("application/json", "text/html"),
			binding.Bind(orgstatus.OrganisationStatus{}),
			authentication.JWTAuth(jwtcfg, PlatformAdminRole),
			orgstatus.Set)

		// Defines /organisations/*/plans resources
		//
		AddDefaultResource(r, fmt.Sprintf("/%v/organisations/:p1/plans", prefixPath), "p2", jwtcfg,
			orgplans.List, orgplans.Get, nil, orgplans.Create, nil,
			orgplans.OrganisationPlanAssoc{})

		// Defines /plans/* resources
		//
		AddDefaultResource(r, fmt.Sprintf("/%v/plans", prefixPath), "p1", jwtcfg,
			plans.List, plans.Get, plans.Delete, plans.Create, plans.Modify,
			plans.Plan{})

		// Defines /roles/* resources
		//
		AddDefaultResource(r, fmt.Sprintf("/%v/roles", prefixPath), "p1", jwtcfg,
			roles.List, roles.Get, roles.Delete, roles.Create, roles.Modify,
			roles.Role{})
	}

	//r.NotFound(strict.MethodNotAllowed, strict.NotFound)
	// m.Action(r.Handle)

	// return m

	return &Server{Name: "Server", Martini: m, Router: r, JWTConfig: jwtcfg}

}
