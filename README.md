![Unide](https://raw.github.com/hivetech/hivetech.github.io/master/images/logo-unide.png)

---------------------------------------------------------------

Hivy
====

**Hivy** is yet an other RESTful interface between http requests and jobs. But it
comes with its modular and simple way to do that and ease the building
of such a popular and efficient web interface.

The goal is to provide a powerful framework that let you focus on your API
design and the available methods your users can reach, without worrying about
so-standard process between.

The project makes heavy use of [etcd](http://coreos.com/docs/etcd/) as a
highly-available centralized configuration storage. Server and clients
preferences are easily stored and accessed from a (soon...) ssl-secured etcd
cluster database.

Out of the box, and because it is part of the **Unide project**, Hivy exposes
awesome [juju](https://juju.ubuntu.com/) commands to authentified remote users.
**Powerful IT infrasctructure building accessible from robust but simple http requests !**


Status
------

[![Build Status](https://drone.io/github.com/hivetech/hivy/status.png)](https://drone.io/github.com/hivetech/hivy/latest)
[![Coverage Status](https://coveralls.io/repos/hivetech/hivy/badge.png?branch=develop)](https://coveralls.io/r/hivetech/hivy?branch=develop)

Branch   | Version
-------- | -----
Stable   | 0.1.5
Develop  | 0.1.6

**Attention** Project is in an *early alpha*, and under heavy development.

Note also I use it to improve my go and devops skills so
there might be some extra dependencies I am testing.


Batteries inluded
-----------------

* Deadly easy to use, build serious applications in minutes
* Http standard authentification
* Permission per user per Method
* Plug and play endpoints, authentification and permission methods
* But not only, up to 2 filters actually before endpoint processing
* Secured, higly-available and centralized configuration storage
* Mysql-ready for users logins
* Debug client provided, dead easy to write one
* Ready for load-balancing and multi-hosts
* Built-in profiling
* 21st century tests

Suit up
-------

First make sure [etcd binary](https://github.com/coreos/etcd/releases/) is available in your $PATH.

```
go get -v github.com/hivetech/hivy
```

Or For development (it will setup etcd)

```console
$ git clone https://github.com/hivetech/hivy.go
$ make
$ make tests
```

Usage
-----

```console
$ make init  # Create admin user and set default hivy configuration
$ ./hivy --help
$ ./hivy -d node -n master --verbose  
$ # Or 
$ forego start

$ # In another terminal
$ curl --user admin:root http://127.0.0.1:8080/v0/actions/user?user=name&pass=pass&group=admin -X PUT
$ curl --user name:pass http://127.0.0.1:8080/v0/actions/dummy  # Test your installation
$ # With the provided clients
$ ./scripts/request v0/actions/help
$ ./scripts/request v0/actions/login?user={user}&pass={pass}
$ ./scripts/request v0/actions/juju/deploy?project={project}&debug=true

$ # Configuration management
$ ./scripts/config set hivy/security/{user}/password secret
$ ./scripts/config get {user}/{project}/services
$ ./scripts/config set {user}/{project}/{charm}/series precise
$ ./scripts/config set {user}/{project}/{charm}/expose True
```

Let's add a new service to our app:

* Implement in endpoints package a method with this signature: ``func (e
  *Endpoint) YourMethod(request *restful.Request, response *restful.Response)``
* Create an authority with authentification and permission
  methods: ``a := NewAuthority(authMethod, permissionMethod)``
* Finally, map the service: ``a.Map("METHOD /path/with/{parameter}", endpoint.YourMethod)``

It is possible to register at a same path multiple endpoints to multiple http
methods. Check out hivy.go to take a glance at it.  We can see as well where to
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


Current API
-----------

Here are listed currently supported methods. With the ./hivy application, all
need user:pass authentification and permissions:

```console
# Admin action methods
PUT /v0/actions/user?user={user}&pass={pass}&group={group}
DELETE /user?user={user}
# User action methods
GET /v0/actions/dummy/
GET /v0/actions/help?method={method}  # method is optionnal
GET /v0/actions/login
GET /v0/actions/juju/status
GET /v0/actions/juju/deploy?project={project}

# Configuration methods
#TODO This is on a different port (4001) for now
GET /v1/keys/{ topology below }
```


Etcd configuration storage topology
-----------------------------------

Etcd storage follows filesystem convention.

```
http://127.0.0.1:4001/v1/keys/hivy/setting1
                                   ...
                                   settingN
                                   security/user1
                                            ...
                                            userN/password
                                                  methods/method1
                                                          ...
                                                          methodN
                                                  ressources/ressource1
                                                             ...
                                                             ressourceN
                             /user1
                             ...
                             /userN/project1
                                    ...
                                   /projectN/setting1
                                             ...
                                             settingN
                                             cell1
                                             ...
                                             cellN/setting1
                                                   ...
                                                   settingN
```

Documentation
-------------

Check it out on [gowalker](http://gowalker.org/github.com/hivetech/hivy),
[godoc](http://godoc.org/github.com/hivetech/hivy), or browse it locally:

```console
$ make doc
$ firefox http://localhost:6060/pkg/github.com/hivetech/hivy/
```


Contributing
------------

> Fork, implement, add tests, pull request, get my everlasting thanks and a
> respectable place here [=)](https://github.com/jondot/groundcontrol)


License
-------

Hivy is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).


---------------------------------------------------------------

![Gophaer](https://raw.github.com/hivetech/hivetech.github.io/master/images/pilotgopher.jpg)
