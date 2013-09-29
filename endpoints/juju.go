package endpoints

import (
    "fmt"
    "path/filepath"
    "strings"
    //"os/exec"

    "github.com/coreos/go-etcd/etcd"
	"github.com/emicklei/go-restful"
)


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
    response, err := storage.Get(filepath.Join("hivy", "charmstore"))
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
        //FIXME It does not have to be local
        //      1. If service contains github url, eventually download it and set charmstore := github_charmstore
        //      2. If not, serach in local
        //      3. Finally try online
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
