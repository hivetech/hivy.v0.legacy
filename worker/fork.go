package main

import (
  "os/exec"
  "fmt"
  "reflect"

  "github.com/benmanns/goworker"
)

func init() {
  hw := newHivyWorker()
  goworker.Register(hw.channel, hw.forker)
}

func (hw *hivyWorker) forker(queue string, args ...interface{}) error {
  log.Infof("trigger %s::forker(%v)\n", queue, args)
  processArgs := []string{}

  process, ok := args[0].(string)
  if !ok { return fmt.Errorf("bad parameter: %v", args[0]) }

  // Converts interface array into string one
  cmdArgsTmp := reflect.ValueOf(args[1])
  for i := 0; i < cmdArgsTmp.Len(); i++ {
    processArgs = append(processArgs, cmdArgsTmp.Index(i).Interface().(string))
  }

  log.Infof("fork %s %v\n", process, processArgs)
  cmd := exec.Command(process, processArgs...)
  output, err := cmd.CombinedOutput()
  if err != nil { 
    log.Errorf("error forking process: %v\n", err)
    return err 
  } 
  log.Infof("successfully forked %s: %v\n", args, string(output))

  return nil
}
