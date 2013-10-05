![Unide](https://raw.github.com/hivetech/hivy/develop/logo-unide.png)

Hivy
====

Hivy is yet an other RESTful interface between http requests and jobs. But it
comes with it's own, modular and simple way to do that and ease the building
of such a popular and efficient web interface.

The goal is to provide a powerful framework that let you focus on your API
design and the available methods your users can reach, without worrying about
so-standard process between.

The project makes heavy use of [etcd](http://coreos.com/docs/etcd/) as a
highly-available centralized configuration storage. Server and clients
preferences are easily stored and accessed from a (soon...) ssl-secured etcd
cluster database.

Out of the box, and because it is part of the Unide project, Hivy exposes
awesome [juju](https://juju.ubuntu.com/) commands to authentified remote users.
Powerful IT infrasctructure building accessible from robust but simple http requests !


Status
------

[![Build Status](https://drone.io/github.com/hivetech/hivy/status.png)](https://drone.io/github.com/hivetech/hivy/latest)
[![Coverage Status](https://coveralls.io/repos/hivetech/hivy/badge.png?branch=develop)](https://coveralls.io/r/hivetech/hivy?branch=develop)

**Attention** Project is in an *early alpha*, and under heavy development.


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
* 21st century tests

Suit-up
-------

First make sure [etcd binary](https://github.com/coreos/etcd/releases/) is available in your $PATH.

```console
$ git clone https://github.com/hivetech/hivy.go
$ make
# That's done, let's check and prepare for further digging
$ make init
$ make tests
$ ./hivy --help
```

Usage
-----

```console
$ make init  # Create admin user and set default hivy configuration
$ ./hivy -d node -n master --verbose
$ # In another terminal
$ curl --user admin:root http://127.0.0.1:8080/createuser?user=name&pass=pass&group=admin
$ curl --user name:pass http://127.0.0.1:8080/dummy  # Test your installation
$ # With the provided client
$ ./client/pencil login
$ ./client/pencil configure --app quantrade --config client/sample-hivy.yml
$ ./client/pencil up --app quantrade
```

To add a new service to your app:

* Implement in endpoints package a method with this signature: ``func (e
  *Endpoint) YourMethod(request *restful.Request, response *restful.Response)``
* In your main(), create an authority with authentification and permission
  methods: ``a := NewAuthority(authMethod, permissionMethod)``
* Finally, still in your main(), register your service: ``a.RegisterGET("/path/with/{parameter}", endpoint.YourMethod)``

You can see as well where to insert custom authentification, permission or any
filter method. Those must be defined in the filters directory and have the
following signature: 

```go
func YourFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
    // Filter whatever you want here
    chain.ProcessFilter(request, response)
}
```


Current API
-----------

Here are listed currently supported methods. With the ./hivy application, all
need user:pass authentification and permissions:

```console
# Admin action methods
GET /createuser?user={user}&pass={pass}&group={group}
# User action methods
GET /dummy/
GET /help?method={method}  # method is optionnal
GET /login?debug={debug}
GET /juju/status?debug={debug}
GET /juju/deploy?project={project}&debug={debug}

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
