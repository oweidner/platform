//  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
// ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
// ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
// ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
// ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
//  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝
//
// Copyright 2015 Codewerft UG (http://www.codewerft.net).
// All rights reserved.

package responses

import (
	"net/http"

	"github.com/gavv/martini-render"
)

type Status struct {
	Status  string
	Message string
}

type GetResponse struct {
	Status
	Results interface{}
}

// XGetError sends an error message back to the client.
func GetError(r render.Render, message string) {
	result := GetResponse{Status{
		Status: "Error", Message: message}, nil}
	r.JSON(http.StatusBadRequest, result)
}

// XGetNotFound sends a 404 message back to the client.
func GetNotFound(r render.Render) {
	status := Status{
		Status: "Error", Message: "Resource not found"}
	r.JSON(http.StatusNotFound, status)
}

// XGetOK sends the result back to the client.
func GetOK(r render.Render, data interface{}) {
	result := GetResponse{Status{
		Status: "OK", Message: ""}, data}
	r.JSON(http.StatusOK, result)
}

// XCreateError sends an error message back to the client.
func CreateError(r render.Render, message string) {
	result := GetResponse{Status{
		Status: "Error", Message: message}, nil}
	r.JSON(http.StatusBadRequest, result)
}

// XCreateOK sends the result back to the client.
func CreateOK(r render.Render, data interface{}) {
	result := GetResponse{Status{
		Status: "OK", Message: ""}, data}
	r.JSON(http.StatusOK, result)
}

// XModifyError sends an error message back to the client.
func ModifyError(r render.Render, message string) {
	result := GetResponse{Status{
		Status: "Error", Message: message}, nil}
	r.JSON(http.StatusBadRequest, result)
}

// XModifyOK sends the result back to the client.
func ModifyOK(r render.Render, data interface{}) {
	result := GetResponse{Status{
		Status: "OK", Message: ""}, data}
	r.JSON(http.StatusOK, result)
}

// XGetError sends an error message back to the client.
func DeleteError(r render.Render, message string) {
	result := GetResponse{Status{
		Status: "Error", Message: message}, nil}
	r.JSON(http.StatusBadRequest, result)
}

// XGetOK sends the result back to the client.
func DeleteOK(r render.Render, message string) {
	result := GetResponse{Status{
		Status: "OK", Message: message}, nil}
	r.JSON(http.StatusOK, result)
}
