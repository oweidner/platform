//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package organizations

import (
	"net/http"

	"codewerft.net/ohoi/apiserver/storage"
	"github.com/gavv/martini-render"
)

// CreateOrganizationRequest is the object that is expected by the
// CreateOrganization() function
type CreateOrganizationRequest struct {
	Name           string  `binding:"required"`
	LocationLat    float32 `binding:"required"`
	LocationLon    float32 `binding:"required"`
	ContactStreet  string  `binding:"required"`
	ContactCity    string  `binding:"required"`
	ContactZip     string  `binding:"required"`
	ContactCountry string  `binding:"required"`
	ContactPhone   string  `binding:"required"`
	ContactEmail   string  `binding:"required"`
}

// CreateOrganization creates a new Organization object in
// the database.
func CreateOrganization(r render.Render, db storage.Datastore, data CreateOrganizationRequest) {

	// stmt, err := db.Get().Prepare(`
	// 	INSERT organization SET name=?,location_lat=?,location_lon=?,
	// 	contact_street=?, contact_city=?, contact_zip=?, contact_country=?, contact_phone=?, contact_email=?`)
	// if err != nil {
	// 	// Forward the error to the client.
	// 	result := responses.CreateError{
	// 		Status{Status: "Error", Message: err.Error()}, []Organization{}}
	// 	r.JSON(http.StatusBadRequest, result)
	// }
	//
	// res, er := stmt.Exec(
	// 	data.Name, data.LocationLat, data.LocationLon, data.ContactStreet,
	// 	data.ContactCity, data.ContactZip, data.ContactCountry,
	// 	data.ContactPhone, data.ContactEmail)
	// if er != nil {
	// 	// Forward the error to the client.
	// 	result := CreateOrganizationResponse{
	// 		Status{Status: "Error", Message: err.Error()}, []Organization{}}
	// 	r.JSON(http.StatusBadRequest, result)
	// }
	// // The id of the newly generated organization
	// organizationId, _ := res.LastInsertId()
	// // Retrieve the newly created object from the database and return it
	// organization, err := DBGetOrganization(db.Get(), organizationId)
	// if er != nil {
	// 	result := CreateOrganizationResponse{
	// 		Status{Status: "Error", Message: err.Error()}, []Organization{}}
	// 	r.JSON(http.StatusOK, result)
	// }
	//
	// result := CreateOrganizationResponse{Status{Status: "OK", Message: "Organization created"},
	// 	[]Organization{organization}}

	r.JSON(http.StatusOK, "")
}
