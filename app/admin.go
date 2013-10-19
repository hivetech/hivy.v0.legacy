package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/emicklei/go-restful"
	"launchpad.net/loggo"

    "github.com/hivetech/hivy"
)

// DeleteUser removes from etcd storage evrything related to the given user-id
func DeleteUser(request *restful.Request, response *restful.Response) {
	user := request.QueryParameter("user")

  var debug bool
  if log.LogLevel() <= loggo.DEBUG { debug = true }
  c := hivy.NewController(user, debug)

	//FIXME Will it delete a directory ?
	feedback, err := c.Delete(filepath.Join("hivy/security", user, "password"))
	if err != nil {
		hivy.HTTPInternalError(response, err)
		return
	}
	log.Debugf("%v\n", feedback)

	//FIXME Ability to delete directory
	//feedback, err = c.Delete(user)
	//if err != nil {
	//hivy.HTTPInternalError(response, err)
	//return
	//}
	//log.Debugf("%v\n", feedback)

	response.WriteEntity(JSON(`{"delete": 0}`))
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
		hivy.HTTPBadRequestError(response, fmt.Errorf("user or pass not provided"))
		return
	}

  var debug bool
  if log.LogLevel() <= loggo.DEBUG { debug = true }
  c := hivy.NewController(user, debug)

	feedback, err := c.Set(filepath.Join("hivy/security", user, "password"), pass, 0)
	if err != nil {
		hivy.HTTPInternalError(response, err)
		return
	}
	log.Debugf("%v\n", feedback)

	feedback, err = c.Set(filepath.Join("hivy/security", user, "ressources/machines"), "0", 0)
	if err != nil {
		hivy.HTTPInternalError(response, err)
		return
	}
	log.Debugf("%v\n", feedback)

    //FIXME v0 hardcoded
	basicAllowedMethods := []string{
		"GET/v0/methods/login",
		"GET/v0/methods/dummy",
		"GET/v0/methods/juju/status", "GET/v0/methods/juju/deploy",
		"GET/v0/methods/help",
	}
	adminAllowedMethods := []string{
		"PUT/v0/methods/user",
		"DELETE/v0/methods/user",
		"GET/v0/methods/juju/bootstrap",
	}

	allowedMethods := basicAllowedMethods
	if group == "admin" {
		allowedMethods = append(allowedMethods, adminAllowedMethods...)
	}

	for _, method := range allowedMethods {
		feedback, err = c.Set(filepath.Join("hivy/security", user, "methods", method), hivy.Allowed, 0)
		if err != nil {
			hivy.HTTPInternalError(response, err)
			return
		}
	}

	response.WriteEntity(JSON(`{"create": 0}`))
}

// Help provides a json object describing available commands
func Help(request *restful.Request, response *restful.Response) {
	method := request.QueryParameter("method")
	json := JSON(fmt.Sprintf(`{"time": "%s"}`, time.Now()))
	//json := EmptyJSON()
	if method == "juju" {
		json.Set("title", "Hivy Juju API")
		json.Set("body", JujuHelp)
		response.WriteEntity(json)
	} else if method == "login" {
		json.Set("title", "Hivy Login API")
		json.Set("body", LoginHelp)
		response.WriteEntity(json)
	} else if method == "user" {
		json.Set("title", "Hivy User API")
		json.Set("body", UserHelp)
		response.WriteEntity(json)
	} else if method == "config" {
		json.Set("title", "Hivy Configuration API")
		json.Set("body", ConfigHelp)
		response.WriteEntity(json)
	} else if method == "help" {
		json.Set("title", "Hivy Help API")
		json.Set("body", HelpHelp)
		response.WriteEntity(json)
	} else {
		json.Set("title", "Hivy API")
		json.Set("resume", GlobalHelp)
		json.Set("/v0/methods/dummy", "Useless so essential, for tests purpose")
		json.Set("/v0/methods/login", "Fetch back an SSL certificate")
		json.Set("/v0/methods/help/{command}", "Get this message, or more details on {command}")
		json.Set("/v0/methods/juju/{command}/{project}", "Manage project through juju commands")
		json.Set("/v1/keys/{path/to/key}", "Set, get, delete settings")
		json.Set("/v0/methods/user?user={user}&pass={pass}&group={group}", "Manage user")
		response.WriteEntity(json)
	}
	return
}
