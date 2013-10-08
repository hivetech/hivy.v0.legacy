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
    "net/http"

	"github.com/emicklei/go-restful"
    "github.com/bitly/go-simplejson"
)


var log = loggo.GetLogger("hivy.endpoints")


//TODO Specifi help should be in specific files ?
const (
    Allowed string = "1"
    Forbidden string = "0"
    Global_help string = `
RESTful framework for busy people. Effectively expose secured, fully configurable background jobs.
`
    Juju_help string = `
Format: GET /juju/{command}?project={project}
With command one of bootstrap, status or deploy method.
The later needs the project parameter as it will read {project} sepcific
configuration and setup accordingly your private cells.  Port exposure and
charms relationships are automatically processed.
`
    Login_help string = `
Format: GET /login
If well authentified, hivy returns a certificate for further secured interactions.
`
    User_help string = `
Format: GET /createuser?user={user}&pass={pass}
Store a new user and its credentials, allowing him to access the rest of the API, restricted to his method permissions.
`
    Config_help string = `
Format: GET /help?method={method}
Will return an help message on the method if provided, global otherwise.
`
)


func HTTPError(writer *restful.Response, err error, httpStatus int) {
    log.Errorf("[HttpError] %v\n", err)
    writer.WriteError(httpStatus, err)
}


func HTTPInternalError(writer *restful.Response, err error) {
    HTTPError(writer, err, http.StatusInternalServerError)
}


func HTTPBadRequestError(writer *restful.Response, err error) {
    HTTPError(writer, err, http.StatusBadRequest)
}


func Json(data string) *simplejson.Json {
    json, err := simplejson.NewJson([]byte(data))
    if err != nil { panic(err) }
    return json
}


func EmptyJSON() *simplejson.Json {
    return Json("{}")
}


/*
 * Note: 
 * Route{method:   restful.WebService.GET,
 *       path:     "/path/{with}/{parameters},
 *       endpoint: Dummy
 * }
 */


// hello-world endpoint, for demo and test purpose
func Dummy(request *restful.Request, response *restful.Response) {
     response.WriteEntity(Json(`{"you": "dummy"}`))
}
