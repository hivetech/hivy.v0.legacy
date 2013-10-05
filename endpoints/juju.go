package endpoints


import (
    "fmt"
    "path/filepath"
    "strings"
    "net/http"
    "os/exec"
    "time"

    "github.com/bitly/go-simplejson"
    "github.com/coreos/go-etcd/etcd"
	"github.com/emicklei/go-restful"

    "github.com/hivetech/hivy/security"
)


const (
    juju_bin string = "juju"
)


func bootstrap(juju string) (*simplejson.Json, error) {
    //FIXME Need sudo permission
    return Json(`{"error": "not implemented"}`), nil
}


func status(juju, user string) (*simplejson.Json, error) {
    //TODO Filter for given user, this will return full system juju status
    log.Infof("Asking for juju status\n")
    cmd := exec.Command("juju", "status",
                        "--format", "json")
    output, err := cmd.CombinedOutput(); 
    if err != nil {
       return EmptyJSON(), err
    }
    json_output, err := simplejson.NewJson(output)
    if err != nil {
       return EmptyJSON(), err
    }
    return json_output, err
}


func deploy(juju string, db *etcd.Client, user string, project string) (*simplejson.Json, error) {
    report   := Json(fmt.Sprintf(`{"time": "%s"}`, time.Now()))
    logs     := []string{}
    deployed := []string{}
    relation := []string{}
    exposed  := []string{}

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
        result, err = db.Get(filepath.Join(user, project, service, "series"))
        if err != nil || len(result) != 1 {
            return EmptyJSON(), err
        }
        series := result[0].Value

        // Deductions
        //FIXME It does not have to be the online cs
        //      1. If service contains github url, eventually download it and set charmstore := github_charmstore
        //      2. If not, serach in local
        //      3. Finally try online
        charm  := fmt.Sprintf("cs:%s/%s", series, service)
        name   := fmt.Sprintf("%s-%s-%s", user, project, service)
        
        // Charm deployment
        //TODO Use CombinedOutput to return logs
        deployed = append(deployed, name)
        log.Infof("Deploying %s (%s) from %s", charm, name, charmstore)
        
        /*
         *cmd := exec.Command("juju",
         *                    "deploy",
         *                    //TODO Detect local or remote repo
         *                    "--repository=" + charmstore,
         *                    charm,
         *                    name)
         *if output, err := cmd.CombinedOutput(); err != nil {
         *   return EmptyJSON(), err
         *} else {
         *   logs = append(logs, string(output))
         *}
         */
         
        result, err = db.Get(filepath.Join(user, project, service, "expose"))
        if err == nil && len(result) == 1 {
            log.Debugf("%v\n", result)
            if result[0].Value == "True" {
                log.Infof("Exposing %s (%s)", charm, name)
                exposed = append(exposed, name)
                
                /*
                 *cmd := exec.Command("juju",
                 *                    "expose",
                 *                    name)
                 *if output, err := cmd.CombinedOutput(); err != nil {
                 *   return EmptyJSON(), err
                 *} else {
                 *   logs = append(logs, string(output))
                 *}
                 */
            }
        }
    }

    // Deployment done, check for relations 
    for _, service := range services {
        name := fmt.Sprintf("%s-%s-%s", user, project, service)
        result, err = db.Get(filepath.Join(user, project, service, "relation"))
        if err == nil && len(result) == 1 {
            log.Debugf("%v\n", result)
            //FIXME needs name or service ?
            relation_target := fmt.Sprintf("%s-%s-%s", user, project, result[0].Value)
            log.Infof("Adding relation between %s and %s", name, relation_target)
            relation = append(relation, fmt.Sprintf("%s->%s", name, relation_target))
            
            /*
             *cmd := exec.Command("juju",
             *                    "add-relation",
             *                    name, relation_target)
             *if output, err := cmd.CombinedOutput(); err != nil {
             *   return EmptyJSON(), err
             *} else {
             *   logs = append(logs, string(output))
             *}
             */
        }
    }

    report.Set("deployed", deployed)
    report.Set("exposed", exposed)
    report.Set("linked", relation)
    report.Set("logs", logs)
    return report, nil
}


// Juju endpoints. It executes the given command, for the given
// user, regarding given project preferences stored in etcd, which are:
//  * charms path if local
//  * charms to be deployed
//  * For each:
//      * Based machine image
func (e *Endpoint) Juju(request *restful.Request, response *restful.Response) {
    //TODO: status and bootstrap does not need project, so this should be a query parameter
    // Parameters
    // Context parameters
    user, _, err := security.Credentials(request)
    if err != nil {
        log.Errorf("[Juju] %v\n", err)
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    command := request.PathParameter("command")

    // Check if juju is available for use
    juju_path, err := exec.LookPath(juju_bin)
	if err != nil {
		log.Criticalf("[boostrap] Unable to find juju program (%s)", juju_bin)
        response.WriteError(http.StatusInternalServerError, err)
        return
	}
	log.Debugf("[bootstrap] juju program available at %s\n", juju_path)

    if request.QueryParameter("debug") == "true" {
        etcd.OpenDebug()
        defer etcd.CloseDebug()
    }
    database := etcd.NewClient()

    if command == "deploy" {
        project := request.QueryParameter("project")
        report, err := deploy(juju_path, database, user, project)
        if err != nil {
            log.Errorf("[Juju] %v\n", err)
            response.WriteError(http.StatusInternalServerError, err)
        } else {
            response.WriteEntity(report)
        }
        return
    } else if command == "status" {
        report, err := status(juju_path, user)
        if err != nil {
            log.Errorf("[Juju] %v\n", err)
            response.WriteError(http.StatusInternalServerError, err)
        } else {
            response.WriteEntity(report)
        }
        return
    } else if command == "bootstrap" {
        report, err := bootstrap(juju_path)
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
