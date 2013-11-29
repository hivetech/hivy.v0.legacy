RESTFul framework
=================

Yet an other RESTful interface between http requests and jobs. But it
comes with its modular and simple way to do that and ease the building
of such a popular and efficient web interface.

The goal is to provide a powerful framework that let you focus on your API
design and the available methods your users can reach, without worrying about
so-standard process between.

The project makes heavy use of [etcd](http://coreos.com/docs/etcd/) as a
highly-available centralized configuration storage. Server and clients
preferences are easily stored and accessed from a (soon...) ssl-secured etcd
cluster database. Please notice it is implemented behind a `Controller`
abstraction that will allow to easily integrate other storage systems in the
futur.


Batteries inluded
-----------------

* Dead easy to use, build serious applications in minutes
* Http standard authentification
* Permission per user per Method
* Plug and play endpoints, authentification and permission methods
* But not only, up to 2 filters actually before endpoint processing
* Secured, higly-available and centralized configuration storage
* Debug client provided, easy to write one
* Ready for load-balancing and multi-hosts
* Built-in profiling
* 21st century tests

Suit up
-------

First make sure [etcd binary](https://github.com/coreos/etcd/releases/) is available in your `$PATH.`

```
go get -v github.com/hivetech/hivy/beacon
```

Or For development (it will setup etcd)

```console
$ git clone https://github.com/hivetech/hivy.git
$ rake install
$ rake tests:check
```

Usage
-----

Let's add a new service to our app:

* Implement a method with this signature: ``func YourMethod(request
  *restful.Request, response *restful.Response)``
* Create an authority with authentification and permission
  methods: ``a := NewAuthority(authMethod, permissionMethod)``
* Finally, map the service: ``a.Map("METHOD /path/with/{parameter}", endpoint.YourMethod)``

It is possible to register at a same path multiple endpoints to multiple http
methods. Check out app/hivy.go to take a glance at it.  We can see as well where to
insert custom authentification, permission or any filter method. Those must
have the following signature: 

```go
func YourFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
    // Filter whatever you want here
    chain.ProcessFilter(request, response)
}
```

Code Example
------------

```go
func authenticate(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
    user, pass, _ := security.Credentials(request)
    if user != "Chuck" || pass != "Norris" {  
        err := fmt.Errorf("you are not chuck norris")
        endpoints.HTTPAuthroizationError(response, err)
        return 
    }
    chain.ProcessFilter(request, response)
}

func control(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
    method := fmt.Sprintf("%s%s", request.Request.Method, request.Request.URL)
    if strings.Contains(method, "deploy") {
        err := fmt.Errorf("deploy method is not supported")
        endpoints.HTTPBadRequestError(response, err)
        return
    }
    chain.ProcessFilter(request, response)
}

func main() {
    router := NewRouter(authenticate, control, profile)

    router.Map("GET juju/{command}", endpoints.Juju)
    router.Map("GET help/", endpoints.Help)

    var userMap = map[string]restful.RouteFunction{
        "PUT user/": endpoints.CreateUser,
        "DELETE user/": endpoints.DeleteUser,
    }
    router.MultiMap(userMap)

    log.Infof("Hivy interface serving on %s\n", url)
    http.ListenAndServe("127.0.0.1:8080", nil)
}
```
