package endpoints


import (
    "fmt"
    "path/filepath"
    "strings"
    "net/http"
    //"os/exec"

    "github.com/bitly/go-simplejson"
    "github.com/coreos/go-etcd/etcd"
	"github.com/emicklei/go-restful"

    "github.com/hivetech/hivy/security"
)


func deploy(db *etcd.Client, user string, project string) (*simplejson.Json, error) {
    // Global settings
    result, err := db.Get(filepath.Join("hivy", "charmstore"))
    if err != nil || len(result) != 1 {
        return EmptyJSON(), err
    }
    charmstore := result[0].Value

    // User defined settings
    result, err = db.Get(filepath.Join(user, project, "services", ))
    if err != nil || len(result) != 1 {
        return EmptyJSON(), err
    }
    services := strings.Split(result[0].Value, ",")

    for _, service := range services {
        //FIXME Should every parameters has default value, or handle it there ?
        result, err = db.Get(filepath.Join(user, project, service, "series", ))
        if err != nil || len(result) != 1 {
            return EmptyJSON(), err
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
         *   return EmptyJSON(), err
         *}
         */

         //TODO If /user/project/service/relation = quelquechose, juju add-relation service quelquechose
/*
 *    ∞ ➜  hivy git:(feature/juju-endpoints) ✗ curl http://127.0.0.1:4001/v1/keys/xav/quantrade/wordpress/relation
 *{"action":"GET","key":"/xav/quantrade/wordpress/relation","value":"wordpress","index":129}
 */
         //TODO If /user/project/service/expose = true, juju expose service
/*
 *    ∞ ➜  hivy git:(feature/juju-endpoints) ✗ curl http://127.0.0.1:4001/v1/keys/xav/quantrade/wordpress/expose
 *{"action":"GET","key":"/xav/quantrade/wordpress/expose","value":"True","index":129}
 */
    }

    return EmptyJSON(), err
}


// Juju endpoints. It executes the given command, for the given
// user, regarding given project preferences stored in etcd, which are:
//  * charms path if local
//  * charms to be deployed
//  * For each:
//      * Based machine image
func (e *Endpoint) Juju(request *restful.Request, response *restful.Response) {
    // Parameters
    // Context parameters
    user := request.QueryParameter("user")
    user, _, err := security.Credentials(request)
    if err != nil {
        log.Errorf("[Juju] %v\n", err)
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    project := request.PathParameter("project")
    command := request.PathParameter("command")
    if user == "" || project == "" || command == "" {
        log.Errorf("[Juju] User, project or command not provided\n")
        response.WriteError(http.StatusBadRequest, fmt.Errorf("User, project or command not provided"))
        return
    }

    etcd.OpenDebug()
    defer etcd.CloseDebug()
    database := etcd.NewClient()

    if command == "deploy" {
        report, err := deploy(database, user, project)
        if err != nil {
            log.Errorf("[Juju] %v\n", err)
            response.WriteError(http.StatusInternalServerError, err)
        } else {
            response.WriteEntity(report)
        }
        return
    }

    response.WriteEntity(EmptyJSON())
}
