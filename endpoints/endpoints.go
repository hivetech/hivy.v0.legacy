package endpoints


import (
    "launchpad.net/loggo"
	"github.com/emicklei/go-restful"
)


//TODO Get hivy.main log level
var log = loggo.GetLogger("hivy.endpoints")


type Endpoint struct {}


func (e *Endpoint) Dummy(req *restful.Request, resp *restful.Response) {
     resp.WriteEntity("{you: dummy}")
}
