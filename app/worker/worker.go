package main

import (
  "fmt"

  "github.com/benmanns/goworker"
  "launchpad.net/loggo"

  "github.com/hivetech/hivy"
)

var log = loggo.GetLogger("hivy.worker")

type hivyWorker struct {
  channel string
  controller *hivy.Controller
}

func newHivyWorker() *hivyWorker {
  return &hivyWorker{
    channel: "Hivy",
    //TODO Hard coded
    controller: hivy.NewController("worker", false),
  }
}

func init() {
  //TODO Hard coded
  logfile := ""
  verbose := true
  modules := []string{"hivy.worker"}
  hivy.SetupLog(modules, verbose, logfile)
}

func main() {
  if err := goworker.Work(); err != nil {
      fmt.Println("Error:", err)
  }
}
