package main

import (
    "fmt"
    "strings"
    "time"

	"github.com/emicklei/go-restful"
    "github.com/davecheney/profile"
)


// Router is the main functionnal entry point
type Router struct {
    authentification restful.FilterFunction
    control restful.FilterFunction
}


// NewRouter needs one function handling user/pass authetification, and one
// function handling method permission for the user who requested it.
func NewRouter(a , c restful.FilterFunction, profiling bool) *Router {
    // Global hook, processed before any service
    //TODO Save it in etcd ?
    restful.Filter(func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
        if profiling {
            cfg := profile.Config {
                Quiet: false,
                ProfilePath: "./profile/" + FormatMethod(request),
                CPUProfile: true,
                MemProfile: true,
                NoShutdownHook: true, // do not hook SIGINT
            }
            p := profile.Start(&cfg)
            defer p.Stop()
            //defer profile.Start(profile.CPUProfile).Stop()
        }

        log.Infof("[global-filter (logger)] %s %s\n", request.Request.Method, request.Request.URL)
        now := time.Now()
        chain.ProcessFilter(request, response)
        log.Infof("[global-filter (timer)] Request processed in %v\n", time.Now().Sub(now))
    })
    return &Router{
        authentification: a,
        control: c,
    }
}


func parsePath(path string) (string, string, error) {
    // We need to separate root path from parameters path
    // Plus, for the sake of readibility, we want the path to contain at least one "/"
    if ! strings.Contains(path, "/") {
        return "", "", fmt.Errorf("one depth path need to end with '/'")
    }
    //TODO Check for path beginning with "/"
    splittedPath := strings.Split(path, "/")

    // Get back together path parameters
    paramPath := ""
    if splittedPath[1] != "" {
        paramPath = "/" + strings.Join(splittedPath[1:], "/")
    }
    return ("/" + splittedPath[0]), paramPath, nil 
}


// Map lets you design your API. When registered to hivy, user-defined callback
// function are processed when "path" is reached by authentified requests
// Example:
//      router.Map("GET", "hello/{world}", func(req, resp) {fmt.Println("Hello world")})
func (a *Router) Map(httpMethod, path string, callback restful.RouteFunction) error {
    //TODO Map(options ...interface{}) {
    log.Infof("Map %s endpoint (%s)\n", path, httpMethod)

    rootPath, paramPath, err := parsePath(path)
    if err != nil { return err }

    // Instanciate a new route at /{path} that returns json data
    ws := new(restful.WebService)
    //TODO To be coherent with etcd, path sould begin with /v1/methods/
    ws.Path(rootPath).
        Consumes("*/*").
	    Produces(restful.MIME_JSON)

    // Create pipeline get-request, according hooks defined at instanciation time 
    // endpoint request -> authentification -> method permission -> callback
    switch {
    default:
        return fmt.Errorf("invalid http method")
    case httpMethod == "GET":
        ws.Route(ws.GET(paramPath).
                    Filter(a.authentification).
                    Filter(a.control).
                    To(callback))
    case httpMethod == "PUT":
        ws.Route(ws.PUT(paramPath).
                    Filter(a.authentification).
                    Filter(a.control).
                    To(callback))
    case httpMethod == "DELETE":
        ws.Route(ws.DELETE(paramPath).
                    Filter(a.authentification).
                    Filter(a.control).
                    To(callback))
    }

    restful.Add(ws)
    return nil
}
