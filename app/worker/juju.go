package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"

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

  goworker.Register(hw.channel, hw.deploy)
}

func (hw *hivyWorker) deploy(queue string, args ...interface{}) error {
  fmt.Printf("trigger %s::deploy(%v)\n", queue, args)

  //TODO Make an jujuWorker, inhereted from hivyWorker
  // Check if juju is available for use
  juju, err := exec.LookPath("juju")
  if err != nil {
    return fmt.Errorf("juju binary not available in path")
  }
  fmt.Printf("juju program available at %s\n", juju)

  name       := args[0].(string)
  service    := args[1].(string)
  series     := args[2].(string)
  charmstore := args[3].(string)

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

  //TODO 1. Dump {user}/{project}/{service}/... into /path/condif.yml
  //     2. deployArgs = append(deployArgs, "--config")
  //        deployArgs = append(deployArgs, /path/config.yml)

  //cmd := exec.Command(juju, deployArgs...)
  //output, err := cmd.CombinedOutput()
  //log.Infof("juju deploy %s: %v\n", name, output)

  //TODO Check what is going on to tell etcd
  return nil
}
