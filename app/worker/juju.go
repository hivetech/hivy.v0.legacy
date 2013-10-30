package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  "reflect"

  "launchpad.net/loggo"
  "github.com/benmanns/goworker"
)

var log = loggo.GetLogger("hivy.worker")

type hivyWorker struct {
  channel string
  //FIXME Controller needs a lot of information: user, project, ...
  //controller *hivy.Controller
}

func newHivyWorker() *hivyWorker {
  return &hivyWorker{
    channel: "Hivy",
  }
}

func init() {
  hw := newHivyWorker()

  goworker.Register(hw.channel, hw.dispatcher)
}

func postDeploy(juju string, args ...interface{}) error {
  log.Infof("trigger postDeploy\n")

  //FIXME I think this command appends much to early
  relations := reflect.ValueOf(args[0])
  for i := 0; i < relations.Len(); i++ {
    // Does not work: services := strings.Split(relations.Index(i).String(), "->")
    services := strings.Split(relations.Index(i).Interface().(string), "->")
    log.Infof("juju add-relation %s->%s\n", services[0], services[1])
    cmd := exec.Command(juju, "add-relation", services[0], services[1])
    output, err := cmd.CombinedOutput()
    if err != nil { return err } 
    log.Infof("juju add-relation %s->%s: %v\n", services[0], services[1], output)
  }

  exposes := reflect.ValueOf(args[1])
  for i := 0; i < exposes.Len(); i++ {
    service := exposes.Index(i).Interface().(string)
    log.Infof("juju expose %s\n", service)
    cmd := exec.Command(juju, "expose", service)
    output, err := cmd.CombinedOutput()
    if err != nil { return err } 
    log.Infof("juju expose %s: %v\n", service, output)
  }

  //TODO Check what is going on to tell etcd
  return nil
}

func deploy(juju string, args ...interface{}) error {
  log.Infof("trigger deploy\n")

  name, ok      := args[0].(string)
  if !ok { return fmt.Errorf("bad parameter") }
  service, _    := args[1].(string)
  series, _     := args[2].(string)
  charmstore, _ := args[3].(string)

  deployArgs := []string{"deploy", "--show-log"}

  //TODO If service contains github url, eventually download it and set charmstore := github_charmstore
  // Default is local storage
  location := "local"
  if _, err := os.Stat(filepath.Join(charmstore, series, service)); os.IsNotExist(err) {
    log.Infof("%s not available localy, trying online store", service)
    location = "cs"
  } else {
    log.Infof("found localcharm %v", service)
    deployArgs = append(deployArgs, "--repository")
    deployArgs = append(deployArgs, charmstore)
  }

  charm := fmt.Sprintf("%s:%s/%s", location, series, service)
  deployArgs = append(deployArgs, charm)
  deployArgs = append(deployArgs, name)

  log.Infof("deploy %s (%s)", charm, name)
  log.Debugf("juju %v\n", deployArgs)

  //TODO 1. Dump {user}/{project}/{service}/... into /path/config.yml
  deployArgs = append(deployArgs, "--config")
  deployArgs = append(deployArgs, "~/dev/projects/hivetech/cells/conf_example.yml")

  cmd := exec.Command(juju, deployArgs...)
  output, err := cmd.CombinedOutput()
  if err != nil { return err } 
  log.Infof("juju deploy %s: %v\n", name, output)

  //TODO Check what is going on to tell etcd
  return nil
}

func (hw *hivyWorker) dispatcher(queue string, args ...interface{}) error {
  log.Infof("trigger %s::dispatcher(%v)\n", queue, args)

  // Check if juju is available for use
  juju, err := exec.LookPath("juju")
  if err != nil {
    return fmt.Errorf("juju binary not available in path")
  }
  log.Infof("juju program available at %s\n", juju)
  //FIXME It only find /usr/bin/juju which is not the good one (i.e. $GOPATH/bin/juju)
  juju = "juju"

  //TODO Both on the same queue, find another way to dispatch (len(args) ?)
  if queue == "deploy" {
    deploy(juju, args...)
  } else if queue == "postDeploy" {
    //Note Not the same queue so won't wait for deploy to finish ?
    postDeploy(juju, args...)
  } else {
    return fmt.Errorf("unknown queue")
  }

  return nil
}
