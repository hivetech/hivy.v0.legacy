package endpoints


import (
    "fmt"
    "net/http"
    "path/filepath"
    //"time"

    "github.com/coreos/go-etcd/etcd"
	"github.com/emicklei/go-restful"
)


func (e *Endpoint) CreateUser(request *restful.Request, response *restful.Response) {
    user := request.QueryParameter("user")
    pass := request.QueryParameter("pass")
    if user == "" || pass == "" {
        log.Errorf("[Juju] User or pass not provided\n")
        response.WriteError(http.StatusBadRequest, fmt.Errorf("User or pass not provided"))
        return
    }

    etcd.OpenDebug()
    defer etcd.CloseDebug()
    db := etcd.NewClient()

    feedback, err := db.Set(filepath.Join("hivy/security", user, "password"), pass, 0)
    if err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    log.Debugf("%v\n", feedback)

    feedback, err = db.Set(filepath.Join("hivy/security", user, "methods", "GET/login"), Allowed, 0)
    if err != nil {
        response.WriteError(http.StatusInternalServerError, err)
        return
    }
    log.Debugf("%v\n", feedback)

    response.WriteEntity(Json(`{"create": 0}`))
}


func (e *Endpoint) Help(request *restful.Request, response *restful.Response) {
    method := request.PathParameter("method")
    //json := Json(fmt.Sprintf(`{"time": %s}`, time.Now()))
    json := EmptyJSON()
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
