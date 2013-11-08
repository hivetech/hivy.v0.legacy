package main

import (
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"launchpad.net/loggo"
	"github.com/emicklei/go-restful"

    "github.com/hivetech/hivy"
)

var log = loggo.GetLogger("hivy.app")

func serveApp(url string, profile bool) {
  // Optional fitering and profiling
	router := hivy.NewRouter(hivy.BasicAuthenticate, hivy.EtcdControlMethod, profile)

  // The router automatically set before "/v{version}/methods/
	// Login function above will be processed when /login path will be
	// reached by authentified requests
	router.Map("GET dummy/", Dummy)
	router.Map("GET help/", Help)
	router.Map("GET login/", Login)
  //FIXME Temporary
  router.Map("PUT plug/", Plug)

	var nodeMap = map[string]restful.RouteFunction{
    //TODO Put on an existing node could upgrade it (i.e. upgrade-charm)
		"PUT node/":    Deploy,
		"DELETE node/": Destroy,
		"GET node/": Status,
		//"PUT node/plug": Plug,
	}
	router.MultiMap(nodeMap)

	//TODO Below line should be allowed (currently method permission forbids it)
	//router.Map("PUT", "user/{user-id}", CreateUser)
	var userMap = map[string]restful.RouteFunction{
		"PUT user/":    CreateUser,
		"DELETE user/": DeleteUser,
	}
	router.MultiMap(userMap)

	log.Infof("Hivy interface serving on %s\n", url)
	http.ListenAndServe(url, nil)
}

func main() {
	// Command line flags configuration
	app := cli.NewApp()
	app.Name = "hivy"
	app.Usage = "Hive router system"
    version := hivy.StableVersion()
	app.Version = version.String()

	// The 2 firsts are hivy's, the lasts are etcd's
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "verbose", Usage: "verbose mode"},
		cli.StringFlag{Name: "logfile", Value: "", Usage: "If specified, will write log there"},
		cli.StringFlag{Name: "listen", Value: "127.0.0.1:8080", Usage: "url to listen to"},
		cli.StringFlag{Name: "n", Value: "", Usage: "the node name (required)"},
		cli.StringFlag{Name: "C", Value: "", Usage: "Node to join"},
		cli.StringFlag{Name: "c", Value: "127.0.0.1:4001", Usage: "Etcd client ip to listen on"},
		cli.StringFlag{Name: "s", Value: "127.0.0.1:7001", Usage: "Raft server ip to listen on"},
		cli.StringFlag{Name: "d", Value: "default-node", Usage: "the directory to store etcd log and snapshot"},
		cli.BoolFlag{Name: "f", Usage: "force new etcd node configuration if existing is found (WARNING: data loss!)"},
		cli.BoolFlag{Name: "profile", Usage: "Profile requests and etcd perfs"},
	}

	// Main function as defined by the cli package
	app.Action = func(c *cli.Context) {
		// Current logger configuration
    modules := []string{"hivy.app", "hivy.worker"}
		hivy.SetupLog(modules, c.Bool("verbose"), c.String("logfile"))
		defer loggo.RemoveWriter("hivy.main")

		// Setup centralized configuration
		// Need to be a new node in the cluster to be ran
		stop := make(chan bool)
		if c.String("n") != "" {
			go hivy.RunEtcd(stop, c.String("n"), c.String("d"), c.String("c"), c.String("s"),
				c.String("C"), c.Bool("f"), c.Bool("verbose"), c.Bool("profile"))
			defer func() {
				stop <- true
			}()
		} else {
			stop = nil
		}

		// Properly shutdown the server when CTRL-C is received
		// Send true on its given channel
		hivy.CatchInterruption(stop)

		// Map routes and start up Hivy server
		serveApp(c.String("listen"), c.Bool("profile"))
	}

	app.Run(os.Args)
}
