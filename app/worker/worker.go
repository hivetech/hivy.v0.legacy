package main

import (
    "fmt"
    "github.com/benmanns/goworker"

    "github.com/hivetech/hivy"
)

func init() {
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
