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
    "./authority"

	"github.com/emicklei/go-restful"
    _ "github.com/go-sql-driver/mysql"
    "os"
    "io/ioutil"
	"net/http"
    "launchpad.net/loggo"
    "github.com/codegangsta/cli"
)

var log = loggo.GetLogger("hivy.authority")

// Object returned when authentified requests reach /login path
// The Cacrt attribute holds a signed certificate that will allow the user to
// interact with hive services.
type Certificate struct {
    Title, Cacrt, Extra string
}

// Endpoint that delivers the above object if a certificate is found.
// It is used as a callback wen registered with a path at the authority server
func login(req *restful.Request, resp *restful.Response) {
    //TODO Generate a new certificate
    ca_data, err := ioutil.ReadFile("ca.crt.example")
    if err != nil {
        log.Errorf("[login] %v\n", err)
		resp.WriteErrorString(404, "404: Could not read certificate")
        return
    }

    // User defined parameter given with http://.../login/{parameter}. Unused for now
	param := req.PathParameter("parameter")
    log.Debugf("Extra info found: %s\n", param)
    // Return the instanciated certificate object
    resp.WriteEntity(Certificate{Title: "Test", Cacrt: string(ca_data), Extra: param})
}

func main() {
    // Command line flags configuration
    app := cli.NewApp()
    app.Name = "Hivy Login"
    app.Usage = "Hivy authentification system"
    app.Version = "0.1.0"

    app.Flags = []cli.Flag {
        cli.BoolFlag{"verbose", "Verbose mode"},
        cli.StringFlag{"listen", "127.0.0.1:8080", "url to listen to"},
    }

    // Main function as defined by the cli package
    app.Action = func(c *cli.Context) {
        // Current logger configuration
        log_level := "hivy.authority=WARNING"
        if c.Bool("verbose") {
            // User wants it more verbose
            log_level = "hivy.authority=TRACE"
        }
        loggo.ConfigureLoggers(log_level)
        log.Debugf("Logging level:", loggo.LoggerInfo())
        defer loggo.RemoveWriter("judo.hivy.watchers")

        log.Infof("Register login service\n")
        // Login function above will be processed when /login path will be
        // reached by authentified requests
        authority.Register("/login", login)
        log.Infof("Hivy authentification serving on %s\n", c.String("listen"))
        //http.ListenAndServe(":8080", nil)
        http.ListenAndServe(c.String("listen"), nil)
    }

    app.Run(os.Args)
}
