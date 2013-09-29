package endpoints

import (
    "io/ioutil"

	"github.com/emicklei/go-restful"
)

// Object returned when authentified requests reach /login path
// The Cacrt attribute holds a signed certificate that will allow the user to
// interact with hive services.
type Certificate struct {
    Status, Cacrt, Extra string
}

// Endpoint that delivers the above object if a certificate is found.
// It is used as a callback wen registered with a path at the authority server
func (e *Endpoints) Login(req *restful.Request, resp *restful.Response) {
    //TODO Generate a new certificate
    log.Debugf("Providing a new certificate")
    ca_data, err := ioutil.ReadFile("ca.crt.example")
    if err != nil {
        //TODO Return proper custom error structure
        log.Errorf("[login] %v\n", err)
		resp.WriteErrorString(404, "404: Could not read certificate")
        return
    }

    // User defined parameter given with http://.../login/{parameter}. Unused for now
	param := req.PathParameter("parameter")
    log.Debugf("Extra info found: %s\n", param)
    // Return the instanciated certificate object
    //TODO Serve static file instead ?
    resp.WriteEntity(Certificate{Status: "OK", Cacrt: string(ca_data), Extra: param})
}
