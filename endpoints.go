package main

import (
    "fmt"
    "io/ioutil"
    "path/filepath"
    "strings"
    //"os/exec"

    "github.com/coreos/go-etcd/etcd"
	"github.com/emicklei/go-restful"
)

type Endpoints struct {}

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

type Context struct {
    User string
    Project string
}

func (e *Endpoints) Deploy(req *restful.Request, resp *restful.Response) {
    // Parameters
    // Context parameters
    user := req.QueryParameter("user")
    project := req.QueryParameter("project")
    if user == "" || project == "" {
        log.Errorf("[deploy] User or project name not provided\n")
        resp.WriteErrorString(406, "406: Not Acceptable")
        return
    }

    etcd.OpenDebug()
    defer etcd.CloseDebug()
    storage := etcd.NewClient()
    // Global settings
    response, err := storage.Get(filepath.Join("hivy", "charmstore" ))
    if err != nil || len(response) != 1 {
        log.Errorf("[deploy] %v\n", err)
        resp.WriteErrorString(404, "404: Not found")
        return
    }
    charmstore := response[0].Value

    // User defined settings
    response, err = storage.Get(filepath.Join(user, project, "services", ))
    if err != nil || len(response) != 1 {
        log.Errorf("[deploy] %v\n", err)
        resp.WriteErrorString(404, "404: Not found")
        return
    }
    services := strings.Split(response[0].Value, ",")

    for _, service := range services {
        //FIXME Should every parameters has default value, or handle it there ?
        response, err = storage.Get(filepath.Join(user, project, service, "series", ))
        if err != nil || len(response) != 1 {
            log.Errorf("[deploy] %v\n", err)
            resp.WriteErrorString(404, "404: Not found")
            return
        }
        series     := response[0].Value

        // Deductions
        charm      := fmt.Sprintf("local:%s/%s", series, service)
        name       := fmt.Sprintf("%s-%s-%s", user, project, service)
        
        // Main cell deployment
        log.Infof("Deploying %s (%s) from %s", charm, name, charmstore)
        /*
         *cmd := exec.Command("juju",
         *                    "deploy",
         *                    //TODO Detect local or remote repo
         *                    "--repository=" + charmstore,
         *                    charm,
         *                    name)
         *if err := cmd.Run(); err != nil {
         *    log.Errorf("[login] %v\n", err)
         *    resp.WriteErrorString(405, "405: Unable to run juju deploy")
         *    return
         *}
         */
    }
     resp.WriteEntity("{juju: deployed}")
}
