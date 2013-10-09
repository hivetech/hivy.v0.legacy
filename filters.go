package main


import (
    "fmt"
    "strings"

	"github.com/emicklei/go-restful"
    "launchpad.net/loggo"

    "github.com/hivetech/hivy/security"
    "github.com/hivetech/hivy/endpoints"
)


func FormatMethod(request *restful.Request) string {
    method := fmt.Sprintf("%s%s", request.Request.Method, request.Request.URL)
    // Consider GET/juju/deploy/*
    if strings.Contains(method, "deploy") {
        method = fmt.Sprintf("%s/%s", request.Request.Method, "juju/deploy")
    }
    param_less_method := strings.Split(method, "?")[0]
    return param_less_method
}


// EtcdControlMethod is a callback part of the request pipeline. It checks in
// etcd if the received request is allowed for the given user.
func EtcdControlMethod(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
    user, _, _ := security.Credentials(request)

    var debug bool = false
    if log.LogLevel() <= loggo.DEBUG {
        debug = true
    }

    controller := NewController(user, debug)

    if err := controller.Update(FormatMethod(request)); err != nil {
        endpoints.HTTPInternalError(response, err)
        return
    }

    is_allowed, err := controller.CheckMethod(FormatMethod(request))
    if err != nil {
        endpoints.HTTPInternalError(response, err)
        return
    } else if ! is_allowed {
        endpoints.HTTPAuthorizationError(response, fmt.Errorf("method disabled"))
        return
    } 
	chain.ProcessFilter(request, response)
}


// BasicAuthenticate is an intermediate step that will check encoded
// credentials before processing the received request.  This function is
// explicitely used in Register() as a filter in the request pipeline.
func BasicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
    //TODO Instead of clear passwor I could stick with encoded or other crypted solution
    // Use base64 decoding to extract from http header user credentials 
    username, passwd, err := security.Credentials(req)
    if err != nil {
        endpoints.HTTPAuthorizationError(resp, err)
		return
    }
    log.Infof("User %s trying to connect with %s\n", username, passwd)

    debug := false
    //TODO Manage a way to plug whatever datastore you want, wherever it is
    ok, err := security.EtcdCheckCredentials(username, passwd, debug)
    if err != nil {
        endpoints.HTTPInternalError(resp, err)
        return 
    }
    if ! ok {
        endpoints.HTTPAuthorizationError(resp, fmt.Errorf("credentials refused"))
		return
	}
    log.Infof("Authentification granted, processing (%s:%s)", username, passwd)
	chain.ProcessFilter(req, resp)
}
