package hivy

import (
	"net/http"

	"github.com/emicklei/go-restful"

	"github.com/hivetech/hivy/security"
  "github.com/hivetech/hivy/beacon"
)

// Get the certificate filename to return
func certificate() (string, error) {
	//TODO Generate a new certificate
	return "ca.crt.example", nil
}

func sshKey() (string, error) {
  return "/home/xavier/.ssh/id_rsa", nil
}

// Login is an endpoint that delivers a certificate, used later for etcd
// communication permission.  It is used as a callback wen registered with a
// path at the authority server
func Login(request *restful.Request, response *restful.Response) {
	user, _, err := security.Credentials(request)
	if err != nil {
		beacon.HTTPInternalError(response, err)
		return
	}
	log.Debugf("Providing a new ssh key to", user)
	key, _ := sshKey()

	// Return the certificate
	http.ServeFile(response.ResponseWriter, request.Request, key)
}
