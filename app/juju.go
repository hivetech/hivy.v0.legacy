package main

import (
  "fmt"
  "os/exec"
  "path/filepath"
  "strings"
  "time"

  "github.com/ghais/goresque"
  "github.com/bitly/go-simplejson"
  "github.com/emicklei/go-restful"
  "launchpad.net/loggo"

  "github.com/hivetech/hivy/security"
  "github.com/hivetech/hivy"
)

const (
  jujuBin string = "juju"
  redisURL string = "127.0.0.1:6379"
)

func bootstrap(juju string) (*simplejson.Json, error) {
  //FIXME Need sudo permission
  return JSON(`{"error": "not implemented"}`), nil
}

func status(juju, user string) (*simplejson.Json, error) {
  //TODO Filter for given user, this will return full system juju status
  //     juju status <user>-<project>-<charm> for each service at {user}/{project}/services
  log.Infof("fetch juju status\n")
  cmd := exec.Command("juju", "status", "--format", "json")
  output, err := cmd.CombinedOutput()
  if err != nil {
    return EmptyJSON(), err
  }

  jsonOutput, err := simplejson.NewJson(output)
  if err != nil {
    return EmptyJSON(), err
  }
  return jsonOutput, err
}

func deploy(juju string, controller *hivy.Controller, user string, project string) (*simplejson.Json, error) {
  report := JSON(fmt.Sprintf(`{"time": "%s"}`, time.Now()))
  logs := []string{}
  deployed := []string{}
  relation := []string{}
  exposed := []string{}

  // Global settings
  result, err := controller.Get(filepath.Join("hivy", "charmstore"))
  if err != nil || len(result) != 1 {
    return EmptyJSON(), err
  }
  charmstore := result[0].Value

  // User defined settings
  result, err = controller.Get(filepath.Join(user, project, "services"))
  if err != nil || len(result) != 1 {
    return EmptyJSON(), err
  }
  services := strings.Split(result[0].Value, ",")

  for _, service := range services {
    //FIXME Should every parameters has default value, or handle it there ?
    result, err = controller.Get(filepath.Join(user, project, service, "series"))
    if err != nil || len(result) != 1 {
      return EmptyJSON(), err
    }
    series := result[0].Value
    name  := fmt.Sprintf("%s-%s-%s", user, project, service)

    // Charm deployment
    deployed = append(deployed, name)
    class := "Hivy"
    queue := "juju"
    client, err := goresque.Dial(redisURL)
    if err != nil { return EmptyJSON(), err }
    client.Enqueue(class, queue, name, service, series, charmstore)

/*
 *    result, err = controller.Get(filepath.Join(user, project, service, "expose"))
 *    if err == nil && len(result) == 1 {
 *      log.Debugf("%v\n", result[0])
 *      if result[0].Value == "True" {
 *        //log.Infof("expose %s (%s)", charm, name)
 *        exposed = append(exposed, name)
 *
 *        cmd := exec.Command(juju, "expose", name)
 *        output, err := cmd.CombinedOutput()
 *        if err != nil { return EmptyJSON(), err }
 *        logs = append(logs, string(output))
 *      }
 *    }
 */
  }

  // Deployment done, check for relations
  //FIXME I think this command appends much to early
/*
 *  for _, service := range services {
 *    result, err = controller.Get(filepath.Join(user, project, service, "relation"))
 *    if err == nil && len(result) == 1 {
 *      name := fmt.Sprintf("%s-%s-%s", user, project, service)
 *      relationTarget := fmt.Sprintf("%s-%s-%s", user, project, result[0].Value)
 *      log.Infof("Adding relation between %s and %s", name, relationTarget)
 *      relation = append(relation, fmt.Sprintf("%s->%s", name, relationTarget))
 *
 *      cmd := exec.Command(juju, "add-relation",
 *      relationTarget, name)
 *      output, err := cmd.CombinedOutput()
 *      if err != nil { return EmptyJSON(), err } 
 *      logs = append(logs, string(output))
 *    }
 *  }
 */

  report.Set("deployed", deployed)
  report.Set("exposed", exposed)
  report.Set("connected", relation)
  report.Set("logs", logs)
  return report, nil
}

// Juju endpoints. It executes the given command, for the given
// user, regarding given project preferences stored in etcd, which are:
//  * charms path if local
//  * charms to be deployed
//  * For each:
//      * Based machine image
func Juju(request *restful.Request, response *restful.Response) {
  // Parameters
  // Context parameters
  user, _, err := security.Credentials(request)
  if err != nil {
    hivy.HTTPInternalError(response, err)
    return
  }
  command := request.PathParameter("command")

  // Check if juju is available for use
  jujuPath, err := exec.LookPath(jujuBin)
  if err != nil {
    hivy.HTTPInternalError(response, err)
    return
  }
  log.Debugf("[bootstrap] juju program available at %s\n", jujuPath)

  var debug bool
  if log.LogLevel() <= loggo.DEBUG { debug = true }
  c := hivy.NewController(user, debug)

  if command == "deploy" {
    project := request.QueryParameter("project")
    report, err := deploy(jujuPath, c, user, project)
    if err != nil {
      hivy.HTTPInternalError(response, err)
    } else {
      response.WriteEntity(report)
    }
    return
  } else if command == "status" {
    report, err := status(jujuPath, user)
    if err != nil {
      hivy.HTTPInternalError(response, err)
    } else {
      response.WriteEntity(report)
    }
    return
  } else if command == "bootstrap" {
    report, err := bootstrap(jujuPath)
    if err != nil {
      hivy.HTTPInternalError(response, err)
    } else {
      response.WriteEntity(report)
    }
    return
  }

  response.WriteEntity(EmptyJSON())
}
