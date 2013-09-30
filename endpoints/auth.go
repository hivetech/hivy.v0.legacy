package endpoints


import (
    "net/http"

	"github.com/emicklei/go-restful"
)


// Endpoint that delivers the above object if a certificate is found.
// It is used as a callback wen registered with a path at the authority server
func (e *Endpoint) Login(request *restful.Request, response *restful.Response) {
    user := request.PathParameter("user")
    //TODO Generate a new certificate
    log.Debugf("Providing a new certificate to", user)

    http.ServeFile(
		response.ResponseWriter,
		request.Request,
		"ca.crt.example")
}
