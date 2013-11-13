package beacon

import (
	"net/http"

	"github.com/emicklei/go-restful"
)

//TODO Better desgin: https://github.com/dotcloud/docker/blob/master/api.go

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
