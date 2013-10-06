package main

import (
    "os"
	"net/http"
    "launchpad.net/loggo"
    "github.com/codegangsta/cli"

    "github.com/hivetech/hivy/endpoints"
)


var log = loggo.GetLogger("hivy.main")


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

        CatchInterruption(stop)

        log.Infof("Hivy interface serving on %s\n", c.String("listen"))
        http.ListenAndServe(c.String("listen"), nil)
    }

    app.Run(os.Args)
}
