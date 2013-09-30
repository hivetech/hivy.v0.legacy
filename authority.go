// Authority package that takes care of user's https requests authentification
//
// This package builds a RESTful server checking user's credentials, and user
// permissions before processing the given callback at the given url path. The
// Register() function instructs the package this parameters. Login and
// password are provided through standard http mechanism and currently verified
// in etcd database after some base64 decoding.
package main

import (
    "strings"
    "time"

	"github.com/emicklei/go-restful"
)

type Authority struct {
    authentification restful.FilterFunction
    control restful.FilterFunction
}

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

// When registered to the authority server, user-defined callback function
// are processed when "path" is reached by authentified requests
// Example:
//      authority.Register("/hello/{world}", func(req, resp) {fmt.Println("Hello world")})
//TODO Not only GET
func (a *Authority) RegisterGET(path string, callback restful.RouteFunction) {
    splitted_path := strings.Split(path, "/")

    ws := new(restful.WebService)
    ws.Path("/" + splitted_path[0]).
        Consumes("*/*").
	    Produces(restful.MIME_JSON)

    pathParameter := ""
    if splitted_path[1] != "" {
        pathParameter = "/" + strings.Join(splitted_path[1:], "/")
    }
    //FIXME A way to plug any basicAuthenticate
    // Create pipeline pipeline get-request -> authentification -> authority -> callback
    ws.Route(ws.GET(pathParameter).
                Filter(a.authentification).
                Filter(a.control).
                To(callback))
    restful.Add(ws)
}
