<p align="center">
  <img src="https://raw.github.com/hivetech/hivetech.github.io/master/images/logo-unide.png" alt="Unide logo"/>
</p>

Hivy
====

Hivy is the restful interface of the **Unide** project that wires users
commands and backends services.

Mostly it lets you setup and build a container powered by the
[dna](https://github.com/hivetech/dna), optimized as a development workspace
and ready to connect to other services as databases, dashboards...

Out of the box it exposes [juju](https://juju.ubuntu.com/) command, users
authentification, and admin related commands.  **Powerful IT infrasctructure
building accessible from robust but simple http requests !**


Status
------

[![Build Status](https://drone.io/github.com/hivetech/hivy/status.png)](https://drone.io/github.com/hivetech/hivy/latest)
[![Coverage Status](https://coveralls.io/repos/hivetech/hivy/badge.png?branch=develop)](https://coveralls.io/r/hivetech/hivy?branch=develop)

Branch   | Version
-------- | -----
Stable   | 0.1.5
Develop  | 0.2.9

**Attention** Project is in an *early alpha*, and under heavy development.

Note also I use it to improve my go and devops skills so
there might be some extra dependencies I am testing.


Suit up
-------

First make sure [etcd binary](https://github.com/coreos/etcd/releases/) is
available in your $PATH. You will need also a redis server for workers (powered
by resque, go port).

```
go get -v github.com/hivetech/hivy
```

Or For development (it will setup etcd)

```console
$ git clone https://github.com/hivetech/hivy.git
$ cd app && rake

$ rake app:run
```

Usage
-----

```console
$ rake app:init  # Create admin user and set default hivy configuration
$ hivy --help
$ hivy -d .conf -n master --verbose  
$ # Or full deployement
$ cd $HIVY_PATH/hivy && forego start

$ # In another terminal
$ # Create a new standard user
$ curl --user admin:root http://127.0.0.1:8080/v0/methods/user?id=name&pass=pass&group=admin -X PUT
$ curl --user name:pass http://127.0.0.1:8080/v0/methods/dummy  # Test your installation

$ # With the provided clients and boxcars proxy on
$ ./scripts/request get help

$ ./scripts/request put node?id=hivelab
$ ./scripts/request put node?id=mysql
$ ./scripts/request put node/plug?id=hivelab&with=mysql

$ ./scripts/request get login > id_rsa
$ ./scripts/request get node?id=hivelab
$ ssh ubuntu@l$YOUR_HOST -p $SSH_PORT -i id_rsa

$ # Configuration management
$ ./scripts/config set hivy/security/{user}/password secret
$ ./scripts/config set {user}/{node}/series precise
$ ./scripts/config set {user}/{node}/expose True
```


API
---

Here are listed currently supported methods. With the `hivy app`, all
need user:pass authentification, and permissions:

```console
# Admin action methods
PUT /v0/methods/user?id={user}&pass={pass}&group={group}
DELETE /user?id={user}
# User methods
GET /v0/methods/dummy/
GET /v0/methods/help?method={method}  # method is optionnal
GET /v0/methods/login
GET /v0/methods/node
GET /v0/methods/node?id={name}
PUT /v0/methods/node?id={name}
DELETE /v0/methods/node?id={name}
PUT /v0/methods/node/plug?id={name}&with={relation}

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
                             /userN/node1
                                    ...
                                   /nodeN/setting1
                                          ...
                                          settingN
                             /backends/localhost/name1
                                                 ...
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
