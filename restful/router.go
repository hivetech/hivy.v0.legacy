package hivy

import (
	"fmt"
	"strings"
	"time"

	"github.com/davecheney/profile"
	"github.com/emicklei/go-restful"
)

// Router is the main functionnal entry point
type Router struct {
    prefix string
	authentification restful.FilterFunction
	control          restful.FilterFunction
	profiling        bool
}

func parsePath(path string) (string, string, error) {
	// We need to separate root path from parameters path
	// Plus, for the sake of readibility, we want the path to contain at least one "/"
	if !strings.Contains(path, "/") {
		return "", "", fmt.Errorf("one-depth paths need to end with '/'")
	}

	var rootPathIndex int
	var paramPath string
	splittedPath := strings.Split(path, "/")

	if splittedPath[rootPathIndex] == "" {
		rootPathIndex++
	}

	// Get back together path parameters
	paramPathIndex := rootPathIndex + 1
	if len(splittedPath) > paramPathIndex {
		if splittedPath[paramPathIndex] != "" {
			paramPath = "/" + strings.Join(splittedPath[paramPathIndex:], "/")
		}
	}
	return (splittedPath[rootPathIndex]), paramPath, nil
}

func parseRequest(request string) (string, string, error) {
	// http method and url path are space separated
	splittedRequest := strings.Split(request, " ")
	if len(splittedRequest) != 2 {
		return "", "", fmt.Errorf("bad request syntax (%s)", request)
	}
	return splittedRequest[0], splittedRequest[1], nil
}

// setupProfiling registers a filter that will measure time to process a
// request alog cpu and memory usage
func setupProfiling() {
	//TODO Save it in etcd ?
	restful.Filter(func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
		cfg := profile.Config{
			Quiet:          false,
			ProfilePath:    "./.profile/" + FormatMethod(request),
			CPUProfile:     true,
			MemProfile:     true,
			NoShutdownHook: true, // do not hook SIGINT
		}
		p := profile.Start(&cfg)
		defer p.Stop()

		log.Infof("[global-filter (logger)] Request received: %s %s\n", request.Request.Method, request.Request.URL)
		now := time.Now()
		chain.ProcessFilter(request, response)
		log.Infof("[global-filter (timer)] Request processed in %v\n", time.Now().Sub(now))
	})
}

// NewRouter needs one function handling user/pass authetification, and one
// function handling method permission for the user who requested it.
func NewRouter(a, c restful.FilterFunction, profiling bool) *Router {
	// Activate etcd, cpu, memory and timer profiling
	if profiling {
		setupProfiling()
	}
	// Hooks processed before any endpoint
	if a == nil {
		a = IdentityFilter
	}
	if c == nil {
		c = IdentityFilter
	}

  version := StableVersion()

	return &Router{
        prefix: fmt.Sprintf("/v%d/methods/", version.major),
		authentification: a,
		control:          c,
		profiling:        profiling,
	}
}

func (a *Router) registerRoute(ws *restful.WebService, httpMethod, paramPath string,
	callback restful.RouteFunction) error {
	// Create pipeline get-request, according hooks defined at instanciation time
	// endpoint request -> authentification -> method permission -> callback
	switch {
	default:
		return fmt.Errorf("invalid http method")
	case httpMethod == "GET":
		ws.Route(ws.GET(paramPath).
			Filter(a.authentification).
			Filter(a.control).
			To(callback))
	case httpMethod == "PUT":
		ws.Route(ws.PUT(paramPath).
			Filter(a.authentification).
			Filter(a.control).
			To(callback))
	case httpMethod == "DELETE":
		ws.Route(ws.DELETE(paramPath).
			Filter(a.authentification).
			Filter(a.control).
			To(callback))
	}
	return nil
}

// Map lets you design your API. When registered to hivy, user-defined callback
// function are processed when "path" is reached by authentified requests
// Example:
//      router.Map("GET hello/{world}", func(req, resp) {fmt.Println("Hello world")})
func (a *Router) Map(request string, callback restful.RouteFunction) error {
	//TODO Map(options ...interface{}) {  // Would bring together Map and MultiMap ?
	httpMethod, path, err := parseRequest(request)
	rootPath, paramPath, err := parsePath(path)
	if err != nil {
		return err
	}

	log.Infof("Map %s endpoint (%s)\n", path, httpMethod)

	// Instanciate a new route at /{path} that returns json data
	ws := new(restful.WebService)
	//TODO To be coherent with etcd, path sould begin with /v1/methods/
	ws.Path(a.prefix + rootPath).
		Consumes("*/*").
		Produces(restful.MIME_JSON)

	err = a.registerRoute(ws, httpMethod, paramPath, callback)
	if err != nil {
		return err
	}

	restful.Add(ws)
	return nil
}

func checkMapping(mapping map[string]restful.RouteFunction) (string, error) {
	var i int
	paths := make([]string, len(mapping))
	for request := range mapping {
		_, path, err := parseRequest(request)
		if err != nil {
			return "", nil
		}
		rootPath, _, err := parsePath(path)
		paths[i] = rootPath
		i++
	}
	return allTheSame(paths)
}

// MultiMap allows app to associate different endpoints under the same root
// path with different http methods
func (a *Router) MultiMap(mapping map[string]restful.RouteFunction) error {
	// All the given keys must have a common root path, return it if it's ok
	rootPath, err := checkMapping(mapping)
	if err != nil {
		return err
	}

	// Instanciate a new route at /{path} that returns json data
	ws := new(restful.WebService)
	ws.Path(a.prefix + rootPath).
		Consumes("*/*").
		Produces(restful.MIME_JSON)

	for request, callback := range mapping {
		// Extract common root path
		httpMethod, path, err := parseRequest(request)
		_, paramPath, err := parsePath(path)
		if err != nil {
			return err
		}

		log.Infof("Map %s endpoint (%s)\n", path, httpMethod)
		err = a.registerRoute(ws, httpMethod, paramPath, callback)
		if err != nil {
			return err
		}
	}
	restful.Add(ws)
	return nil
}
