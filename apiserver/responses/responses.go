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

	"github.com/martini-contrib/render"
)

type Status struct {
	Status  string
	Message string
	Details string
}

type Response struct {
	Status
	Length int
	Result interface{}
}

func Error(r render.Render, message string) {
	result := Status{Status: "Error", Message: message}
	r.JSON(http.StatusBadRequest, result)
}

// XGetNotFound sends a 404 message back to the client.
func ResourceNotFound(r render.Render) {
	status := Status{Status: "Error", Message: "Resource not found"}
	r.JSON(http.StatusNotFound, status)
}

// XGetNotFound sends a 404 message back to the client.
func AuthenticationError(r render.Render, message string, details string) {
	status := Status{Status: "Error", Message: message, Details: details}
	r.JSON(http.StatusUnauthorized, status)
}

// XCreateOK sends the result back to the client.
func OKStatusPlusData(r render.Render, data interface{}, length int) {
	result := Response{
		Status: Status{Status: "OK", Message: ""},
		Result: data,
		Length: length}
	r.JSON(http.StatusOK, result)
}

// XCreateOK sends the result back to the client.
func OKStatusOnly(r render.Render, message string) {
	status := Status{Status: "OK", Message: message}
	r.JSON(http.StatusOK, status)
}
