package endpoints


import (
    "fmt"
    "path/filepath"
    "time"

    "launchpad.net/loggo"
    "github.com/coreos/go-etcd/etcd"
	"github.com/emicklei/go-restful"
)


// DeleteUser removes from etcd storage evrything related to the given user-id
func DeleteUser(request *restful.Request, response *restful.Response) {
    user := request.QueryParameter("user")

    if log.LogLevel() <= loggo.DEBUG {
        etcd.OpenDebug()
        defer etcd.CloseDebug()
    }
    db := etcd.NewClient()

    //FIXME Will it delete a directory ?
    feedback, err := db.Delete(filepath.Join("hivy/security", user, "password"))
    if err != nil {
        HTTPInternalError(response, err)
        return
    }
    log.Debugf("%v\n", feedback)

    //FIXME Ability to delete directory
    //feedback, err = db.Delete(user)
    //if err != nil {
        //HTTPInternalError(response, err)
        //return
    //}
    //log.Debugf("%v\n", feedback)

    response.WriteEntity(Json(`{"delete": 0}`))
}


// CreateUser Stores given credentials and creates methods permission for the
// given user, regarding the given group (only admin and standard group are
// supported at the moment)
func CreateUser(request *restful.Request, response *restful.Response) {
    user := request.QueryParameter("user")
    pass := request.QueryParameter("pass")
    //TODO group specific permissions
    group := request.QueryParameter("group")
    if user == "" || pass == "" {
        HTTPBadRequestError(response, fmt.Errorf("user or pass not provided")) 
        return
    }
    if group == "" {
        group = "basic"
    }

    if log.LogLevel() <= loggo.DEBUG {
        etcd.OpenDebug()
        defer etcd.CloseDebug()
    }
    db := etcd.NewClient()

    feedback, err := db.Set(filepath.Join("hivy/security", user, "password"), pass, 0)
    if err != nil {
        HTTPInternalError(response, err)
        return
    }
    log.Debugf("%v\n", feedback)

    feedback, err = db.Set(filepath.Join("hivy/security", user, "ressources/machines"), "0", 0)
    if err != nil {
        HTTPInternalError(response, err)
        return
    }
    log.Debugf("%v\n", feedback)

    basicAllowedMethods := []string {
        "GET/login",
        "GET/dummy",
        "GET/juju/status", "GET/juju/deploy",
        "GET/help",
    }
    adminAllowedMethods := []string {
        "PUT/user",
        "DELETE/user",
        "GET/juju/bootstrap",
    }

    allowedMethods := basicAllowedMethods
    if group == "admin" {
        allowedMethods = append(allowedMethods, adminAllowedMethods...)
    }

    for _, method := range allowedMethods {
        feedback, err = db.Set(filepath.Join("hivy/security", user, "methods", method), Allowed, 0)
        if err != nil {
            HTTPInternalError(response, err)
            return
        }
    }

    response.WriteEntity(Json(`{"create": 0}`))
}

// Help provides a json object describing available commands
func Help(request *restful.Request, response *restful.Response) {
    method := request.QueryParameter("method")
    json := Json(fmt.Sprintf(`{"time": "%s"}`, time.Now()))
    //json := EmptyJSON()
    if method == "juju" {
        json.Set("title", "Hivy Juju API")
        json.Set("body", Juju_help)
        response.WriteEntity(json)
    } else if method == "login" {
        json.Set("title", "Hivy Login API")
        json.Set("body", Login_help)
        response.WriteEntity(json)
    } else if method == "user" {
        json.Set("title", "Hivy User API")
        json.Set("body", User_help)
        response.WriteEntity(json)
    } else if method == "config" {
        json.Set("title", "Hivy Configuration API")
        json.Set("body", Config_help)
        response.WriteEntity(json)
    } else {
        json.Set("title", "Hivy API")
        json.Set("resume", Global_help)
        json.Set("/dummy", "Useless so essential, for tests purpose")
        json.Set("/login", "Fetch back an SSL certificate")
        json.Set("/help/{command}", "Get this message, or more details on {command}")
        json.Set("/juju/{command}/{project}", "Manage project through juju commands")
        json.Set("/v1/keys/{path/to/key}", "Set, get, delete settings")
        json.Set("/createuseruser={user}&pass={pass}", "Store new user credentials")
        response.WriteEntity(json)
    }
    return
}
