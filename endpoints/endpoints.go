// Package endpoints is a collection of functions that can be registered to the
// go-restful package as routes.
//
// After authentification and permission, the user request is processed by
// those functions that take care to return informations to him as well.
//
// Usage example:
//      authority.Map("GET dummy/", endpoint.Dummy)
package endpoints

import (
	"launchpad.net/loggo"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var log = loggo.GetLogger("hivy.endpoints")

const (
	// Allowed is a macro representing an accessible method
	Allowed string = "true"
	// Forbidden is a macro representing an hiden method
	Forbidden string = "false"
)

// httpFactoryError logs the error and writes back a standard http message
func httpFactoryError(writer *restful.Response, err error, httpStatus int) {
	log.Errorf("[httpfactoryerror] %v\n", err)
	writer.WriteError(httpStatus, err)
}

// HTTPInternalError handles server errors
func HTTPInternalError(writer *restful.Response, err error) {
	httpFactoryError(writer, err, http.StatusInternalServerError)
}

// HTTPBadRequestError handles unknown requests
func HTTPBadRequestError(writer *restful.Response, err error) {
	httpFactoryError(writer, err, http.StatusBadRequest)
}

// HTTPAuthorizationError handles permission failure
func HTTPAuthorizationError(writer *restful.Response, err error) {
	writer.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
	httpFactoryError(writer, err, http.StatusUnauthorized)
}

// JSON converts string to simple.Json object
func JSON(data string) *simplejson.Json {
	json, err := simplejson.NewJson([]byte(data))
	if err != nil {
		panic(err)
	}
	return json
}

// EmptyJSON is a commodity shortcut
func EmptyJSON() *simplejson.Json {
	return JSON("{}")
}

/*
 * Note:
 * Route{method:   restful.WebService.GET,
 *       path:     "/path/{with}/{parameters},
 *       endpoint: Dummy
 * }
 */

// Dummy is the hello-world endpoint, for demo and test purpose
func Dummy(request *restful.Request, response *restful.Response) {
	response.WriteEntity(JSON(`{"you": "dummy"}`))
}
