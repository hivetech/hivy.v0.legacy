package endpoints


import (
    "fmt"
    "path/filepath"
    "strings"
    "net/http"
    //"os/exec"

    "github.com/coreos/go-etcd/etcd"
	"github.com/emicklei/go-restful"
)


// Juju deploy <charm> endpoint. It deploys the given project for the given
// user, regarding preferences stored in etcd, which are:
//  * charms path if local
//  * charms to be deployed
//  * For each:
//      * Based machine image
func (e *Endpoint) Deploy(request *restful.Request, response *restful.Response) {
    // Parameters
    // Context parameters
    //TODO Get user from header, using security.Credentials()
    user := request.QueryParameter("user")
    project := request.PathParameter("project")
    if user == "" || project == "" {
        log.Errorf("[deploy] User or project name not provided\n")
        response.WriteError(http.StatusBadRequest, fmt.Errorf("User or project name not provided"))
        return
    }

    etcd.OpenDebug()
    defer etcd.CloseDebug()
    storage := etcd.NewClient()
    // Global settings
    result, err := storage.Get(filepath.Join("hivy", "charmstore"))
    if err != nil || len(result) != 1 {
        log.Errorf("[deploy] %v\n", err)
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    charmstore := result[0].Value

    // User defined settings
    result, err = storage.Get(filepath.Join(user, project, "services", ))
    if err != nil || len(result) != 1 {
        log.Errorf("[deploy] %v\n", err)
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    services := strings.Split(result[0].Value, ",")

    for _, service := range services {
        //FIXME Should every parameters has default value, or handle it there ?
        result, err = storage.Get(filepath.Join(user, project, service, "series", ))
        if err != nil || len(result) != 1 {
            log.Errorf("[deploy] %v\n", err)
            response.WriteError(http.StatusInternalServerError, err)
            return
        }
        series     := result[0].Value

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

         //TODO If /user/project/service/relation = quelquechose, juju add-relation service quelquechose
         //TODO If /user/project/service/expose = true, juju expose service
    }
     response.WriteEntity("{juju: deployed}")
}
