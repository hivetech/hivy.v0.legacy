package main

import (
    "fmt"
    "os"
    "os/exec"
    "os/signal"
	"net/http"
    "launchpad.net/loggo"
    "github.com/codegangsta/cli"

    "github.com/hivetech/hivy/endpoints"
)

var log = loggo.GetLogger("hivy.main")

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
func RunEtcd(stop chan bool, name string, directory string, force bool, verbose bool, profiling string) {
    //TODO End it properly: http://blog.labix.org/2011/10/09/death-of-goroutines-under-control
    // etcd command line arguments
    args := []string{"-n", name, "-d", directory, "--cpuprofile", profiling}
    if force {
        args = append(args, "-f")
    }
    if verbose {
        args = append(args, "-v")
    }

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

    // Central log level configuration
    for _, module := range app_modules {
        if err := loggo.ConfigureLoggers(module + "=" + log_level.String()); err != nil {
            return fmt.Errorf("configure loggers")
        }
    }

    log.Debugf("Logging level:", loggo.LoggerInfo())
    return nil
}


//TODO Read git tag to make it automatic
// Modest accessor
func Version() string { return current_version }


func main() {
    // Command line flags configuration
    app := cli.NewApp()
    app.Name = "hivy"
    app.Usage = "Hive router system"
    app.Version = Version()

    // The 2 firsts are hivy's, the lasts are etcd's
    app.Flags = []cli.Flag {
        cli.BoolFlag{"verbose", "verbose mode"},
        cli.StringFlag{"logfile", "", "If specified, will write log there"},
        cli.StringFlag{"listen", "127.0.0.1:8080", "url to listen to"},
        cli.StringFlag{"n", "dafault-name", "the node name (required)"},
        cli.StringFlag{"d", ".", "the directory to store etcd log and snapshot"},
        cli.BoolFlag{"f", "force new etcd node configuration if existing is found (WARNING: data loss!)"},
        cli.StringFlag{"cpuprofile", "", "write cpu profile to file"},
    }

    // Main function as defined by the cli package
    app.Action = func(c *cli.Context) {
        // Current logger configuration
        SetupLog(c.Bool("verbose"), c.String("logfile"))
        defer loggo.RemoveWriter("hivy.main")

        // Setup centralized configuration
        stop := make(chan bool)
        go RunEtcd(stop, c.String("n"), c.String("d"), c.Bool("f"), c.Bool("verbose"), c.String("cpuprofile"))
        defer func() {
            stop <- true
        }()

        //TODO Makes it possible to omit one or all method
        authority := NewAuthority(BasicAuthenticate, EtcdControlMethod)
        // Available application services
        //NOTE This is currently useless
        var endpoint endpoints.Endpoint

        // Login function above will be processed when /login path will be
        // reached by authentified requests
        authority.RegisterGET("login/", endpoint.Login)
        authority.RegisterGET("juju/{command}", endpoint.Juju)
        //FIXME Overflow when "/" is missing
        authority.RegisterGET("dummy/", endpoint.Dummy)
        //TODO Delete user using DELETE http method
        authority.RegisterGET("createuser/", endpoint.CreateUser)
        authority.RegisterGET("help/", endpoint.Help)

        log.Infof("Setup exit method")
        ctrl_c := make(chan os.Signal, 1)
        signal.Notify(ctrl_c, os.Interrupt)
        go func() {
            // Stuck still ctrl-c interruption
            for sig := range ctrl_c {
                log.Infof("[main] Server interrupted (%v), cleaning...", sig)
                // End etcd instance
                stop <- true
                os.Exit(0)
            }
        }()

        log.Infof("Hivy interface serving on %s\n", c.String("listen"))
        http.ListenAndServe(c.String("listen"), nil)
    }

    app.Run(os.Args)
}
