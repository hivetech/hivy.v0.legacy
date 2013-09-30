package filters


import (
    "fmt"
    "net/http"

	"github.com/emicklei/go-restful"

    "../security"
)


// Intermediate step that will check encoded credentials before processing the received request.
// This function is explicitely used in Register() as a filter in the request pipeline.
func BasicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
    //TODO Instead of clear passwor I could stick with encoded or other crypted solution
    username, passwd, err := security.Credentials(req)
    if err != nil {
        log.Errorf("[basicAuthenticate] %v\n", err)
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteError(http.StatusUnauthorized, err)
		return
    }
    log.Infof("User %s trying to connect with %s\n", username, passwd)

    //TODO Manage a way to plug whatever datastore you want, wherever it is
    ok, err := security.EtcdCheckCredentials(username, passwd)
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
