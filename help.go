package hivy

import (
  "fmt"
	"time"

	"github.com/emicklei/go-restful"
)

const (
	// GlobalHelp shows a sum up of the app
	GlobalHelp string = `
RESTful framework for busy people. Effectively expose secured, fully configurable background jobs.
`

	// NodeHelp explains juju exposed endpoints
	NodeHelp string = `
Format: GET v0/methods/node?id={name}
Fetch back informations about machine <name>, if provided, your global
infrastructure otherwise. It shows as well the ssh forwarded port of your
lab.
Format: PUT v0/methods/node?id={name}
Deploy machine <name> according its default setup and the configuration
stored ast v1/keys/user/{name}/...
Format: DELETE v0/methods/node?id={name}
Remove service {name} and its associated container.
Format: PUT v0/methods/node/plug?id={name}&with={relation}
Make {relation} available to use for machine {name}.
`

	// LoginHelp explains what happens when logging in
	LoginHelp string = `
Format: GET /login
If well authentified, hivy returns a private ssh key for further secured interactions.
`

	// UserHelp details endpoints relative to user management
	UserHelp string = `
Format: GET /createuser?id={user}&pass={pass}
Store a new user and its credentials, allowing him to access the rest of the
API, restricted to his method permissions.
Format: DELETE /createuser?id={user}
Remove user permissions from hivy
`

	// HelpHelp explains the help endpoint
	HelpHelp string = `
Format: GET /help/{method}
Will return an help message on the {method} if provided, global otherwise.
Available topics are help, login, config, user, node
`

	// ConfigHelp explains etcd interaction
	ConfigHelp string = `
Format: GET v1/keys/{path}
Interact with the etcd database, with classic set and get methods.
Note: Currently accessible through an other port
`
)

// Help provides a json object describing available commands
func Help(request *restful.Request, response *restful.Response) {
	//method := request.QueryParameter("method")
	method := request.PathParameter("method")
	json := JSON(fmt.Sprintf(`{"time": "%s"}`, time.Now()))
	if method == "node" {
		json.Set("title", "Hivy machines API")
		json.Set("body", NodeHelp)
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
		json.Set("get /v0/methods/dummy", "Useless, for tests purpose")
		json.Set("get /v0/methods/login", "Fetch back your SSH private key")
		json.Set("get /v0/methods/help/{method}", "Get this message, or more details on {command}")
		json.Set("put|delete|get /v0/methods/node?id={name}", "Manage your infrastructure")
		json.Set("put /v0/methods/node/plug?id={name}&with={relation}", "Wire two machines")
		json.Set("get|delete /v0/methods/user?id={user}&pass={pass}&group={group}", "Manage users")
		json.Set("get|put|delete /v1/keys/{path/to/key}", "Set, get, delete settings")
		response.WriteEntity(json)
	}
	return
}
