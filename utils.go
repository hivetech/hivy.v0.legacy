package main

import (
	"fmt"
	"launchpad.net/loggo"
	"os"
	"os/exec"
	"os/signal"
)

const (
	superVerboseLogLevel loggo.Level = loggo.TRACE
	verboseLogLevel      loggo.Level = loggo.INFO
	defaultLogLevel      loggo.Level = loggo.WARNING
	currentVersion       string      = "0.1.0"
	// Change here to absolute or relative path if etcd is not in your $PATH
	etcdBin string = "etcd"
)

// RunEtcd launchs an http-based key-value storage that holds user and system
// configuration. Here is spawned a new instance, restricted to relevant
// command line flags for hivy application.
func RunEtcd(stop chan bool, name, directory, clientIP, raftIP, clusterIP string,
	force, verbose bool, profile bool) {
	//TODO End it properly: http://blog.labix.org/2011/10/09/death-of-goroutines-under-control
	// etcd command line arguments
	args := []string{
		"-n", name,
		"-d", directory,
		"-c", clientIP,
		"-s", raftIP,
	}
	if force {
		args = append(args, "-f")
	}
	if force {
		args = append(args, "-f")
	}
	if profile {
		args = append(args, "--cpuprofile")
		args = append(args, "./profile/etcd-profile")
	}
	if clusterIP != "" {
		args = append(args, "-C")
		args = append(args, clusterIP)
	}

	log.Debugf("%v\n", args)
	etcdPath, err := exec.LookPath(etcdBin)
	if err != nil {
		log.Criticalf("[runetcd] unable to find etcd program")
		return
	}
	log.Debugf("[runetcd] etcd program available at %s\n", etcdPath)

	// Spawn the process
	cmd := exec.Command("etcd", args...)
	if err := cmd.Start(); err != nil {
		log.Errorf("[main.runEtcd] %v\n", err)
		return
	}
	//TODO Get some output ?
	log.Infof("etcd server successfully started")
	// Wait for stop instruction
	<-stop
}

// SetupLog set application's modules log level
func SetupLog(verbose bool, logfile string) error {
	//TODO logfile handler
	var appModules = []string{
		"hivy.main",
		"hivy.endpoints",
		"hivy.security",
	}
	logLevel := defaultLogLevel
	if verbose {
		logLevel = superVerboseLogLevel
	}

	// If a file was specified, replace console output
	if logfile != "" {
		target, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		fileWriter := loggo.NewSimpleWriter(target, &loggo.DefaultFormatter{})
		//loggo.RegisterWriter("logfile", file_writer, log_level)
		_, err = loggo.ReplaceDefaultWriter(fileWriter)
		if err != nil {
			return err
		}
	}

	// Central log level configuration
	for _, module := range appModules {
		if err := loggo.ConfigureLoggers(module + "=" + logLevel.String()); err != nil {
			return err
		}
	}

	log.Debugf("logging level:", loggo.LoggerInfo())
	return nil
}

// Version is a Modest accessor
//TODO Read git tag to make it automatic
//     return fmt.Sprintf("%d.%d.%d", major, minor, fix)
func Version() string { return currentVersion }

// CatchInterruption handles SIGINT signal to clean the application before
// exiting. If th stop channel exists it will trigger a signal usuable elsewhere
func CatchInterruption(stop chan bool) {
	log.Infof("Setup exit method")
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)
	go func() {
		// Stuck still ctrl-c interruption
		for sig := range ctrlC {
			log.Infof("[main] server interrupted (%v), cleaning...", sig)
			// End etcd instance
			if stop != nil {
				stop <- true
			}
			os.Exit(0)
		}
	}()
}

// allTheSame check the equlity of every string values in the array
func allTheSame(values []string) (string, error) {
	for i, v := range values {
		if i == (len(values) - 1) {
			break
		} else if v != values[i+1] {
			return "", fmt.Errorf("different strings in the array")
		}
	}
	// If we arrived here, all values are identical
	return values[0], nil
}
