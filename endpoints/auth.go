package endpoints


import (
    "net/http"

	"github.com/emicklei/go-restful"

    "github.com/hivetech/hivy/security"
)


// Get the certificate filename to return
func certificate() (string, error) {
    //TODO Generate a new certificate
    return "ca.crt.example", nil
}


// Endpoint that delivers a certificate, used later for etcd communication
// permission.  It is used as a callback wen registered with a path at the
// authority server
func (e *Endpoint) Login(request *restful.Request, response *restful.Response) {
    user, _, err := security.Credentials(request)
    if err != nil {
        log.Errorf("[Juju] %v\n", err)
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    log.Debugf("Providing a new certificate to", user)
    cert_file, _ := certificate()

    // Return the certificate
    http.ServeFile(
		response.ResponseWriter,
		request.Request,
		cert_file)
}
