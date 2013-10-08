package main


import (
    "fmt"
    "net/http"
    "strings"

	"github.com/emicklei/go-restful"

    "github.com/hivetech/hivy/security"
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


// Callback part of the request pipeline. It checks in etcd if the received
// request is allowed for the given user.
func EtcdControlMethod(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
    user, _, _ := security.Credentials(request)

    //FIXME debug must not be asked for by the client
    var debug bool = false
    if request.QueryParameter("debug") == "true" {
        debug = true
    }

    controller := NewController(user, debug)

    if err := controller.Update(FormatMethod(request)); err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }

    is_allowed, err := controller.CheckMethod(FormatMethod(request))
    if err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    } else if ! is_allowed {
        response.WriteError(http.StatusUnauthorized, fmt.Errorf("Method disabled"))
        return
    } 
	chain.ProcessFilter(request, response)
}


// Intermediate step that will check encoded credentials before processing the received request.
// This function is explicitely used in Register() as a filter in the request pipeline.
func BasicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
    //TODO Instead of clear passwor I could stick with encoded or other crypted solution
    // Use base64 decoding to extract from http header user credentials 
    username, passwd, err := security.Credentials(req)
    if err != nil {
        log.Errorf("[basicAuthenticate] %v\n", err)
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteError(http.StatusUnauthorized, err)
		return
    }
    log.Infof("User %s trying to connect with %s\n", username, passwd)

    debug := false
    //TODO Manage a way to plug whatever datastore you want, wherever it is
    ok, err := security.EtcdCheckCredentials(username, passwd, debug)
    if err != nil {
        log.Errorf("[basicAuthenticate] %v\n", err)
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteError(http.StatusInternalServerError, err)
        return 
    }
    if ! ok {
        log.Warningf("Authentification failed (%s:%s)\n", username, passwd)
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteError(http.StatusUnauthorized, fmt.Errorf("credentials not accepted"))
		return
	}
    log.Infof("Authentification granted, processing (%s:%s)", username, passwd)
	chain.ProcessFilter(req, resp)
}
