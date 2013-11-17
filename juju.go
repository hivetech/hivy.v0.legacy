package hivy

import (
  "fmt"
  "time"
  "path/filepath"
  "os"
  "os/exec"
  "strings"
  "io/ioutil"

  "launchpad.net/goyaml"
  "github.com/ghais/goresque"
  "launchpad.net/loggo"
  "github.com/bitly/go-simplejson"

  "github.com/hivetech/hivy/beacon"
)

//TODO Get and set juju methods should be handled by the controller.
//     Some etcd watchers on service config to trigger juju set attr=value

const (
  //TODO Series policy (automatic choice ?)
  defaultSeries string = "precise"
  jujuBin  string = "juju"
  workerClass string = "Hivy"
)

// Juju is a provider used for Nodes management
type Juju struct {
  Path string
  Controller *beacon.Controller
}

// NewJuju initializes juju provider informations
func NewJuju() (*Juju, error) {
  // User information is not yet relevant for the controller
  user := ""
  // Check if juju is available for use
  jp, err := exec.LookPath(jujuBin)
  if err != nil { return nil, err }
  log.Debugf("[bootstrap] juju program available at %s\n", jp)

  var debug bool
  if log.LogLevel() <= loggo.DEBUG { debug = true }
  c := beacon.NewController(user, debug)

  return &Juju{
    Path: jp,
    Controller: c, 
  }, nil
}

// id returns the way nodes are called with juju provider
func (jj *Juju) id(user, service string) string {
  return fmt.Sprintf("%s-%s", user, service)
}

// Charmstore search for the appropriate configuration
func (jj *Juju) Charmstore(service string) (string, string, error) {
  //TODO If service contains github url, eventually download it and set charmstore := github_charmstore

  result, err := jj.Controller.Get(filepath.Join("hivy", "charmstore"))
  if err != nil || len(result) != 1 {
    return "", "", err
  }
  path := result[0].Value

  // Default is local storage
  prefix := "local"
  if _, err := os.Stat(filepath.Join(path, defaultSeries, service)); os.IsNotExist(err) {
    log.Infof("%s not available localy, use online store", service)
    prefix = "cs"
  }   
  return path, prefix, nil
}

// fetchConfig uses controller to read a configuration and dump it into a yaml file
func (jj *Juju) fetchConfig(user, service string) (string, error) {
  //FIXME This works only for hivetest and hivelab
  //      Could explore directories instead
  if ! strings.HasPrefix(service, "hive") {
    log.Infof("service not supported for custom configuration")
    return "", nil
  }
  // Initialize
  configLabTree := []string{
    "serf-ip", "dna-repo", "targets", "app-repo", "extra-packages",
    "openlibs", "editor", "terminal_multiplexer", "plugins",
    "shell/default", "shell/prompt",
    "dev/node_version", "dev/go_version", "dev/python_version",
  }
  config := make(map[string]string)

  // Fetch and store configuration
  log.Infof("reading configuration for node %s\n", jj.id(user, service))
  for _, key := range configLabTree {
    result, err := jj.Controller.Get(filepath.Join(user, service, key))
    if err == nil && len(result) == 1 {
      log.Infof("[config] %s = %s\n", strings.Replace(key, "/", ".", -1), result[0].Value)
      config[strings.Replace(key, "/", "-", -1)] = result[0].Value 
    }
  }

  // Dump file
  charmConfig := make(map[string]interface{})
  charmConfig[jj.id(user, service)] = config
  log.Infof("serializing config")
  out, err := goyaml.Marshal(charmConfig)
  if err != nil { return "", err }
  filePath := filepath.Join("/tmp", user + "-config.yaml")
  log.Infof("dump configuration into %s\n", filePath)
  if err := ioutil.WriteFile(filePath, out, 0666); err != nil {
    return "", err
  }

  return filePath, nil
}

// Status fetches given service informations
func (jj *Juju) Status(user, service string) (*simplejson.Json, error) {
  //TODO Without id, should return every services matching {user}-*, not everything
  id := jj.id(user, service)
  args := []string{"status", "--format", "json"}
  if service != "" {
    args = append(args, id)
  }
  log.Infof("fetch juju status (%s)\n", id)

  cmd := exec.Command("juju", args...)
  output, err := cmd.CombinedOutput()
  if err != nil {
    return EmptyJSON(), err
  }
  log.Debugf("successful request: %v\n", string(output))

  jsonOutput, err := simplejson.NewJson(output)
  if err != nil {
    return EmptyJSON(), err
  }

  mapping, _ := vmSshForward(user, jj.Controller, jsonOutput)
  jsonOutput.Set("ssh-port", mapping)

  return jsonOutput, err
}

// Deploy uses juju deploy to create a new service
func (jj *Juju) Deploy(user, service string) (*simplejson.Json, error) {
  args := []string{"deploy", "--show-log"}
  id := jj.id(user, service)
  report := JSON(fmt.Sprintf(`{"time": "%s"}`, time.Now()))
  log.Infof("deploy juju service: %s\n", id)

  // Get charms location
  storePath, storePrefix, err := jj.Charmstore(service)
  if err != nil { return EmptyJSON(), err }
  if storePrefix == "local" {
    args = append(args, "--repository")
    args = append(args, storePath)
  }

  // Add final service syntax to deploy
  args = append(args, fmt.Sprintf("%s:%s/%s", storePrefix, defaultSeries, service))
  args = append(args, id)

  // Read and dump user configuration
  confPath, err := jj.fetchConfig(user, service)
  if err != nil { return EmptyJSON(), err }
  if confPath != "" {
    args = append(args, "--config") 
    args = append(args, confPath) 
  }

  // Charm deployment
  //log.Infof("enqueue process")
  //client, err := goresque.Dial(redisURL)
  //if err != nil { return EmptyJSON(), err }
  //client.Enqueue(workerClass, "fork", jj.Path, args)

  report.Set("deployed", id)
  report.Set("provider", "juju")
  report.Set("arguments", args)
  report.Set("series", defaultSeries)
  return report, nil
}

// Destroy uses juju destroy to remove a service
func (jj *Juju) Destroy(user, service string) (*simplejson.Json, error) {
  id := jj.id(user, service)
  report := JSON(fmt.Sprintf(`{"time": "%s"}`, time.Now()))
  log.Infof("destroy juju service: %s\n", id)

  //Note For now this is massive and basic destruction
  unitArgs := []string{"destroy-unit", id + "/0", "--show-log"}
  serviceArgs := []string{"destroy-service", id, "--show-log"}

  cmd := exec.Command("juju", "status", id, "--format", "json")
  output, err := cmd.CombinedOutput()
  if err != nil { return EmptyJSON(), err }
  status, err := simplejson.NewJson(output)
  machineID, err := status.GetPath("services", id, "units", id+"/0", "machine").String()
  if err != nil { return EmptyJSON(), err }
  machineArgs := []string{"destroy-machine", machineID, "--show-log"}

  client, err := goresque.Dial(redisURL)
  if err != nil { return EmptyJSON(), err }
  log.Infof("enqueue destroy-unit")
  client.Enqueue(workerClass, "fork", jj.Path, unitArgs)
  time.Sleep(5 * time.Second)
  log.Infof("enqueue destroy-service")
  client.Enqueue(workerClass, "fork", jj.Path, serviceArgs)
  time.Sleep(5 * time.Second)
  log.Infof("enqueue destroy-machine")
  client.Enqueue(workerClass, "fork", jj.Path, machineArgs)

  report.Set("provider", "juju")
  report.Set("unit destroyed", id + "/0")
  report.Set("service destroyed", id)
  report.Set("machine destroyed", machineID)
  report.Set("unit arguments", unitArgs)
  report.Set("service arguments", serviceArgs)
  report.Set("machine arguments", machineArgs)
  return report, nil
}

// AddRelation links two juju services
func (jj *Juju) AddRelation(user, serviceOne, serviceTwo string) (*simplejson.Json, error) {
  idOne := jj.id(user, serviceOne)
  idTwo := jj.id(user, serviceTwo)
  report := JSON(fmt.Sprintf(`{"time": "%s"}`, time.Now()))
  log.Infof("add juju relation: %s -> %s\n", idOne, idTwo)

  args := []string{"add-relation", "--show-log", idOne, idTwo}
  client, err := goresque.Dial(redisURL)
  if err != nil { return EmptyJSON(), err }
  log.Infof("enqueue add-relation")
  client.Enqueue(workerClass, "fork", jj.Path, args)

  report.Set("provider", "juju")
  report.Set("plugged", fmt.Sprintf("%s -> %s", idOne, idTwo))
  report.Set("relation arguments", args)
  return report, nil
}
