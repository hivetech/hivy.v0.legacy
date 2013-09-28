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
	"net/http"
    "launchpad.net/loggo"
    "github.com/codegangsta/cli"
)

var log = loggo.GetLogger("hivy.main")

func main() {
    // Command line flags configuration
    app := cli.NewApp()
    app.Name = "hivy"
    app.Usage = "Hivy authentification system"
    app.Version = "0.1.0"

    app.Flags = []cli.Flag {
        cli.BoolFlag{"verbose", "Verbose mode"},
        cli.StringFlag{"listen", "127.0.0.1:8080", "url to listen to"},
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
        log.Debugf("Logging level:", loggo.LoggerInfo())
        defer loggo.RemoveWriter("hivy.main")

        var endpoints Endpoints

        log.Infof("Register login service\n")
        // Login function above will be processed when /login path will be
        // reached by authentified requests
        Register("/login", endpoints.Login)

        log.Infof("Register Deploy service\n")
        Register("/deploy", endpoints.Deploy)

        log.Infof("Hivy interface serving on %s\n", c.String("listen"))
        http.ListenAndServe(c.String("listen"), nil)
    }

    app.Run(os.Args)
}
