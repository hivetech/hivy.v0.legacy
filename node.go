package hivy

import (
  "fmt"
  "path/filepath"

  "github.com/bitly/go-simplejson"
  "github.com/emicklei/go-restful"

  "github.com/hivetech/hivy/beacon"
  "github.com/hivetech/hivy/security"
)

const (
  provider string = "juju"
  redisURL string = "127.0.0.1:6379"
  labNode  string = "hivelab"
)

// vmSshForward returns the host port the service is using for ssh
func vmSshForward(user string, controller *beacon.Controller, status *simplejson.Json) (string, error) {
  // We want user's lab hostname, ok like <user>-local-machine-<ID>
  // And we have everything but the machine ID where the lab is
  //FIXME serviceKey is juju format specific, should use id() method for proper abstraction
  serviceKey := fmt.Sprintf("%s-%s", user, labNode)
  //machineID, err := status.Get("services").Get(serviceKey).Get("units").Get(serviceKey+"/0").Get("machine").String()
  machineID, err := status.GetPath("services", serviceKey, "units", serviceKey+"/0", "machine").String()
  if err != nil {
    return "", err
  }
  log.Debugf("got hivelab machine id: %s", machineID)

  result, err := controller.Get(filepath.Join("hivy", "mapping", "xavier-local-machine-"+machineID))
  if err != nil {
    return "", err
  }
  return result.Value, nil
}

//TODO Factorizes 3 functions that are identical but the juju method

// Status fetchs informations about the given node id
func Status(request *restful.Request, response *restful.Response) {
  user, _, err := security.Credentials(request)
  if err != nil {
    beacon.HTTPInternalError(response, err)
    return
  }

  id := request.QueryParameter("id")

  if provider == "juju" {
    juju, err := NewJuju()
    if err != nil {
      beacon.HTTPInternalError(response, err)
      return
    }
    report, err := juju.Status(user, id)
    if err != nil {
      beacon.HTTPInternalError(response, err)
    } else {
      response.WriteEntity(report)
    }
    return
  }
  beacon.HTTPInternalError(response, err)
}

// Deploy creates new nodes
func Deploy(request *restful.Request, response *restful.Response) {
  user, _, err := security.Credentials(request)
  if err != nil {
    beacon.HTTPInternalError(response, err)
    return
  }

  id := request.QueryParameter("id")

  if provider == "juju" {
    //TODO Deploy command on existing service triggers upgrade-charm
    //TODO For suitable charms, it could also deploy other clustered units
    juju, err := NewJuju()
    if err != nil {
      beacon.HTTPInternalError(response, err)
      return
    }
    report, err := juju.Deploy(user, id)
    if err != nil {
      beacon.HTTPInternalError(response, err)
    } else {
      response.WriteEntity(report)
    }
    return
  }
  beacon.HTTPInternalError(response, err)
}

// Destroy removes nodes
func Destroy(request *restful.Request, response *restful.Response) {
  user, _, err := security.Credentials(request)
  if err != nil {
    beacon.HTTPInternalError(response, err)
    return
  }

  id := request.QueryParameter("id")

  if provider == "juju" {
    juju, err := NewJuju()
    if err != nil {
      beacon.HTTPInternalError(response, err)
      return
    }
    report, err := juju.Destroy(user, id)
    if err != nil {
      beacon.HTTPInternalError(response, err)
    } else {
      response.WriteEntity(report)
    }
    return
  }
  beacon.HTTPInternalError(response, err)
}

// Plug allows interactions between two nodes
func Plug(request *restful.Request, response *restful.Response) {
  user, _, err := security.Credentials(request)
  if err != nil {
    beacon.HTTPInternalError(response, err)
    return
  }

  id := request.QueryParameter("id")
  withID := request.QueryParameter("with")

  if provider == "juju" {
    juju, err := NewJuju()
    if err != nil {
      beacon.HTTPInternalError(response, err)
      return
    }
    report, err := juju.AddRelation(user, id, withID)
    if err != nil {
      beacon.HTTPInternalError(response, err)
    } else {
      response.WriteEntity(report)
    }
    return
  }
  beacon.HTTPInternalError(response, err)
}
