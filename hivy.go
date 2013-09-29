// Application using authority package to deliver certificates to authentified
// users (i.e. users that provided with their request a login and password
// registered in mysql 'hive' database).
// Protocol used here is REST.
//
// Usage example:
//      $ go run restful-basic-authentification --verbose --listen 0.0.0.0:8080
// Client usage example
//      $ GET http://localhost:8080/login
//      $ curl --user name:password http://127.0.0.1:8080/login
//      $ http --auth admin:admin http://127.0.0.1:8080/login/deathstar
//      $ python -m "import requests; requests.get('http://127.0.0.1:8080/login/deathstar', auth=('user', 'passwd')'"
package main

import (
    "os"
    "os/exec"
	"net/http"
    "launchpad.net/loggo"
    "github.com/codegangsta/cli"

    "./endpoints"
)

var log = loggo.GetLogger("hivy.main")

func runEtcd(stop chan bool, name string, directory string, force bool, verbose bool, profiling string) {
    args := []string{"-n", name, "-d", directory, "--cpuprofile", profiling}
    if force {
        args = append(args, "-f")
    }
    if verbose {
        args = append(args, "-v")
    }

    cmd := exec.Command("etcd", args...)
    if err := cmd.Start(); err != nil {
        log.Errorf("[main.runEtcd] %v\n", err)
        return
    }
    //TODO Get some output ?
    log.Infof("Etcd server successfully started")
    <- stop
}

func main() {
    // Command line flags configuration
    app := cli.NewApp()
    app.Name = "hivy"
    app.Usage = "Hivy authentification system"
    app.Version = "0.1.0"

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
        log_level := "WARNING"
        if c.Bool("verbose") {
            // User wants it more verbose
            log_level = "TRACE"
        }
        loggo.ConfigureLoggers("hivy.main=" + log_level)
        log.Debugf("Main logging level:", loggo.LoggerInfo())
        defer loggo.RemoveWriter("hivy.main")

        // Setup centralized configuration
        stop := make(chan bool)
        go runEtcd(stop, c.String("n"), c.String("d"), c.Bool("f"), c.Bool("verbose"), c.String("cpuprofile"))
        defer func() {stop <- true}()

        // Available application services
        var endpoints endpoints.Endpoints
        log.Infof("Register Login endpoint\n")
        // Login function above will be processed when /login path will be
        // reached by authentified requests
        Register("/login", endpoints.Login)

        log.Infof("Register Deploy endpoint\n")
        Register("/deploy", endpoints.Deploy)

        log.Infof("Register Dummy endpoint\n")
        Register("/dummy", endpoints.Dummy)

        log.Infof("Hivy interface serving on %s\n", c.String("listen"))
        http.ListenAndServe(c.String("listen"), nil)
    }

    app.Run(os.Args)
}
