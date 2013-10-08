package main

import (
    "fmt"
    "os"
    "os/exec"
    "os/signal"
    "launchpad.net/loggo"
)

const (
    verbose_log_level loggo.Level = loggo.TRACE
    default_log_level loggo.Level = loggo.WARNING
    current_version string = "0.1.0"
    // Change here to absolute or relative path if etcd is not in your $PATH
    etcd_bin string = "etcd"
)

// Etcd is an http-based key-value storage that holds user and system
// configuration. Here is spawned a new instance, restricted to relevant
// command line flags for hivy application.
func RunEtcd(stop chan bool, name, directory, client_ip, raft_ip, cluster_ip string, 
             force, verbose bool, profile bool) {
    //TODO End it properly: http://blog.labix.org/2011/10/09/death-of-goroutines-under-control
    // etcd command line arguments
    args := []string{
        "-n", name, 
        "-d", directory, 
        "-c", client_ip, 
        "-s", raft_ip, 
    }
    if force { args = append(args, "-f") }
    if force { args = append(args, "-f") }
    if profile { 
        args = append(args, "--cpuprofile") 
        args = append(args, "./profile/etcd-profile") 
    }
    if cluster_ip != "" { 
        args = append(args, "-C") 
        args = append(args, cluster_ip) 
    }

    log.Warningf("%v\n", args)
    etcd_path, err := exec.LookPath(etcd_bin)
	if err != nil {
		log.Criticalf("[RunEtcd] Unable to find etcd program")
        return
	}
	log.Debugf("[RunEtcd] etcd program available at %s\n", etcd_path)

    // Spawn the process
    cmd := exec.Command("etcd", args...)
    if err := cmd.Start(); err != nil {
        log.Errorf("[main.runEtcd] %v\n", err)
        return
    }
    //TODO Get some output ?
    log.Infof("Etcd server successfully started")
    // Wait for stop instruction
    <- stop
}


// Set application's modules log level
func SetupLog(verbose bool, logfile string) error {
    //TODO logfile handler
    var app_modules = []string{
        "hivy.main",
        "hivy.endpoints",
        "hivy.security",
    }
    log_level := default_log_level
    if verbose {log_level = verbose_log_level}

    // If a file was specified, replace console output
    if logfile != "" {
        target, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil { return err }
		file_writer := loggo.NewSimpleWriter(target, &loggo.DefaultFormatter{})
        //loggo.RegisterWriter("logfile", file_writer, log_level)
		_, err = loggo.ReplaceDefaultWriter(file_writer)
		if err != nil { return err }
    }

    // Central log level configuration
    for _, module := range app_modules {
        if err := loggo.ConfigureLoggers(module + "=" + log_level.String()); err != nil {
            return err
        }
    }

    log.Debugf("Logging level:", loggo.LoggerInfo())
    return nil
}


//TODO Read git tag to make it automatic
// Modest accessor
func Version() string { return current_version }


func CatchInterruption(stop chan bool) {
    log.Infof("Setup exit method")
    ctrl_c := make(chan os.Signal, 1)
    signal.Notify(ctrl_c, os.Interrupt)
    go func() {
        // Stuck still ctrl-c interruption
        for sig := range ctrl_c {
            log.Infof("[main] Server interrupted (%v), cleaning...", sig)
            // End etcd instance
            if stop != nil { stop <- true }
            os.Exit(0)
        }
    }()
}


func allTheSame(values []string) (string, error) {
    for i, v := range values {
        if i == (len(values) - 1) {
            break
        } else if v != values[i + 1] {
            return "", fmt.Errorf("different strings in the array")
        }
    }
    // If we arrived here, all values are identical
    return values[0], nil
}
