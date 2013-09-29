package endpoints

import (
    "launchpad.net/loggo"
	"github.com/emicklei/go-restful"
)

//TODO Get hivy.main log level
var log = loggo.GetLogger("hivy.endpoints")

type Endpoints struct {}

func init() {
    log_level := "TRACE"
    loggo.ConfigureLoggers("hivy.endpoints=" + log_level)
}

func (e *Endpoints) Dummy(req *restful.Request, resp *restful.Response) {
     resp.WriteEntity("{you: dummy}")
}
