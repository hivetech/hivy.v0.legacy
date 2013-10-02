// Authority takes care of user's https requests permission
//
// It checks user's credentials, and user permissions before processing the
// given callback at the given url path. The Register() function instructs the
// package this parameters. Login and password are provided through standard
// http mechanism and currently verified in etcd database after some base64
// decoding.
package main

import (
    "strings"
    "time"

	"github.com/emicklei/go-restful"
)

// Main entry point
type Authority struct {
    authentification restful.FilterFunction
    control restful.FilterFunction
}

// Constructor. Needs one function handling user/pass authetification, and one
// function handling method permission for the user who requested it.
func NewAuthority(a restful.FilterFunction, c restful.FilterFunction) *Authority {
    // Global hook, processed before any service
    restful.Filter(func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
        log.Infof("[global-filter (logger)] %s %s\n", req.Request.Method, req.Request.URL)
        now := time.Now()
        chain.ProcessFilter(req, resp)
        log.Infof("[global-filter (timer)] Request processed in %v\n", time.Now().Sub(now))
    })
    return &Authority{
        authentification: a,
        control: c,
    }
}

// When registered to hivy, user-defined callback function
// are processed when "path" is reached by authentified requests
// Example:
//      authority.RegisterGET("/hello/{world}", func(req, resp) {fmt.Println("Hello world")})
func (a *Authority) RegisterGET(path string, callback restful.RouteFunction) {
    //TODO Not only GET

    // We need to separate root path from parameters path
    splitted_path := strings.Split(path, "/")

    // Instanciate a new route at /{path} that returns json data
    ws := new(restful.WebService)
    ws.Path("/" + splitted_path[0]).
        Consumes("*/*").
	    Produces(restful.MIME_JSON)

    // Get back together path parameters
    pathParameter := ""
    if splitted_path[1] != "" {
        pathParameter = "/" + strings.Join(splitted_path[1:], "/")
    }
    // Create pipeline get-request, according hooks defined at instanciation time 
    // endpoint request -> authentification -> method permission -> callback
    ws.Route(ws.GET(pathParameter).
                Filter(a.authentification).
                Filter(a.control).
                To(callback))
    restful.Add(ws)
}
