// Router between user requests and Hive jobs. The client sends authentified
// http requests to reach endpoints defined in the endpoints directory.
// Through the process a centralized configuration server, backed by etcd,
// stores user-defined and hive settings, and a central authority is asked for
// methods permissions.
// Communication protocol used here is REST.
//
// Usage example:
//      $ go run hivy --verbose --listen 0.0.0.0:8080
// Client usage example
//      $ GET -C name:pass http://localhost:8080/login/name
//      $ curl --user name:pass http://127.0.0.1:8080/login/name
//      $ python -m "import requests; requests.get('http://127.0.0.1:8080/login/name', auth=('user', 'pass')'"
package main

import (
    "os"
    "os/exec"
    "os/signal"
	"net/http"
    "launchpad.net/loggo"
    "github.com/codegangsta/cli"

    "github.com/hivetech/hivy/endpoints"
    "github.com/hivetech/hivy/filters"
)

var log = loggo.GetLogger("hivy.main")

// Etcd is an http-based key-value storage that holds user and system
// configuration. Here is spawned a new instance, restricted to relevant
// command line flags for hivy application.
func RunEtcd(stop chan bool, name string, directory string, force bool, verbose bool, profiling string) {
    // etcd command line arguments
    args := []string{"-n", name, "-d", directory, "--cpuprofile", profiling}
    if force {
        args = append(args, "-f")
    }
    if verbose {
        args = append(args, "-v")
    }

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
    //TODO Should end etcd process ?
}


// Set application's modules log level
func SetupLog(verbose bool) {
    app_modules := []string{
        "hivy.main",
        "hivy.endpoints",
        "hivy.filters",
        "hivy.security",
    }

    log_level := "WARNING"
    if verbose {log_level = "TRACE"}

    // Central log level configuration
    for _, module := range app_modules {
        loggo.ConfigureLoggers(module + "=" + log_level)
    }
    log.Debugf("Logging level:", loggo.LoggerInfo())
}


func Version() string { return "0.1.0" }


func main() {
    // Command line flags configuration
    app := cli.NewApp()
    app.Name = "hivy"
    app.Usage = "Hive router system"
    app.Version = Version()

    // The 2 firsts are hivy's, the lasts are etcd's
    app.Flags = []cli.Flag {
        cli.BoolFlag{"verbose", "verbose mode"},
        cli.StringFlag{"listen", "127.0.0.1:8080", "url to listen to"},
        cli.StringFlag{"n", "dafault-name", "the node name (required)"},
        cli.StringFlag{"d", ".", "the directory to store etcd log and snapshot"},
        cli.BoolFlag{"f", "force new etcd node configuration if existing is found (WARNING: data loss!)"},
        cli.StringFlag{"cpuprofile", "", "write cpu profile to file"},
    }

    // Main function as defined by the cli package
    app.Action = func(c *cli.Context) {
        // Current logger configuration
        SetupLog(c.Bool("verbose"))
        defer loggo.RemoveWriter("hivy.main")

        // Setup centralized configuration
        stop := make(chan bool)
        go RunEtcd(stop, c.String("n"), c.String("d"), c.Bool("f"), c.Bool("verbose"), c.String("cpuprofile"))
        defer func() {
            stop <- true
        }()

        //TODO Makes it possible to omit one or all method
        authority := NewAuthority(filters.BasicAuthenticate, filters.EtcdControl)
        // Available application services
        //NOTE This is currently useless
        var endpoint endpoints.Endpoint

        log.Infof("Register Login endpoint\n")
        // Login function above will be processed when /login path will be
        // reached by authentified requests
        authority.RegisterGET("login/{user}", endpoint.Login)

        log.Infof("Register Deploy endpoint\n")
        authority.RegisterGET("deploy/{project}", endpoint.Deploy)

        //FIXME Overflow when "/" is missing
        log.Infof("Register Dummy endpoint\n")
        authority.RegisterGET("dummy/", endpoint.Dummy)

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
