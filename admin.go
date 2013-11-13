package hivy

import (
	"fmt"
	"path/filepath"

	"github.com/emicklei/go-restful"
	"launchpad.net/loggo"

    "github.com/hivetech/hivy/beacon"
)

// DeleteUser removes from etcd storage evrything related to the given user-id
func DeleteUser(request *restful.Request, response *restful.Response) {
	user := request.QueryParameter("id")

  var debug bool
  if log.LogLevel() <= loggo.DEBUG { debug = true }
  c := beacon.NewController(user, debug)

	feedback, err := c.Delete(filepath.Join("hivy/security", user, "password"))
	if err != nil {
		beacon.HTTPInternalError(response, err)
		return
	}
	log.Debugf("%v\n", feedback)

	//FIXME Ability to delete directory
	//feedback, err = c.Delete(user)
	//if err != nil {
	//beacon.HTTPInternalError(response, err)
	//return
	//}
	//log.Debugf("%v\n", feedback)

	response.WriteEntity(JSON(`{"deleted": 1}`))
}

// CreateUser Stores given credentials and creates methods permission for the
// given user, regarding the given group (only admin and standard group are
// supported at the moment)
func CreateUser(request *restful.Request, response *restful.Response) {
	user := request.QueryParameter("id")
	pass := request.QueryParameter("pass")
	//TODO group specific permissions
	group := request.QueryParameter("group")
	if user == "" || pass == "" {
		beacon.HTTPBadRequestError(response, fmt.Errorf("user or pass not provided"))
		return
	}

  var debug bool
  if log.LogLevel() <= loggo.DEBUG { debug = true }
  c := beacon.NewController(user, debug)

	feedback, err := c.Set(filepath.Join("hivy/security", user, "password"), pass, 0)
	if err != nil {
		beacon.HTTPInternalError(response, err)
		return
	}
	log.Debugf("%v\n", feedback)

	feedback, err = c.Set(filepath.Join("hivy/security", user, "ressources/machines"), "0", 0)
	if err != nil {
		beacon.HTTPInternalError(response, err)
		return
	}
	log.Debugf("%v\n", feedback)

    //FIXME v0 hardcoded
	basicAllowedMethods := []string{
		"GET/v0/methods/login",
		"GET/v0/methods/dummy",
		"PUT/v0/methods/node/plug", "GET/v0/methods/node", "DELETE/v0/methods/node",
		//"GET/v0/methods/help",
    //TODO Find a way to allow subsequent paths, like help/*
		"GET/v0/methods/help/node",
		"GET/v0/methods/help/config",
		"GET/v0/methods/help/login",
		"GET/v0/methods/help/user",
		"GET/v0/methods/help/help",
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
		feedback, err = c.Set(filepath.Join("hivy/security", user, "methods", method), beacon.Allowed, 0)
		if err != nil {
			beacon.HTTPInternalError(response, err)
			return
		}
	}

	response.WriteEntity(JSON(`{"created": 1}`))
}
