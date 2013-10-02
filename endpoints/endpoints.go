// Collection of functions that can be registered to the go-restful package as
// routes.
//
// After authentification and permission, the user request is processed by
// those functions that take care to return informations to him as well.
//
// Usage example:
//      authority.RegisterGET("dummy/", endpoint.Dummy)
package endpoints


import (
    "launchpad.net/loggo"
	"github.com/emicklei/go-restful"
)


var log = loggo.GetLogger("hivy.endpoints")


type Endpoint struct {}


// hello-world endpoint, for demo and test purpose
func (e *Endpoint) Dummy(request *restful.Request, response *restful.Response) {
     response.WriteEntity("{you: dummy}")
}
