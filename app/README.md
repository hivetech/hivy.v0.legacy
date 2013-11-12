Hivy application
================

Suit up
-------

First make sure [etcd binary](https://github.com/coreos/etcd/releases/) is
available in your $PATH. You will need also a redis server for workers (powered
by resque, go port).

```
go get -v github.com/hivetech/hivy/app
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
$ hivy -d node -n master --verbose  
$ # Or full deployement
$ forego start

$ # In another terminal
$ curl --user admin:root http://127.0.0.1:8080/v0/methods/user?user=name&pass=pass&group=admin -X PUT
$ curl --user name:pass http://127.0.0.1:8080/v0/methods/dummy  # Test your installation
$ # With the provided clients
$ ./scripts/request get v0/methods/help
$ ./scripts/request get v0/methods/login
$ ./scripts/request put v0/methods/node?id={project}

$ # Configuration management
$ ./scripts/config set hivy/security/{user}/password secret
$ ./scripts/config set {user}/{node}/series precise
$ ./scripts/config set {user}/{node}/expose True
```


Current API
-----------

Here are listed currently supported methods. With the `hivy app`, all
need user:pass authentification, and permissions:

```console
# Admin action methods
PUT /v0/methods/user?user={user}&pass={pass}&group={group}
DELETE /user?user={user}
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


Contributing
------------

> Fork, implement, add tests, pull request, get my everlasting thanks and a
> respectable place here [=)](https://github.com/jondot/groundcontrol)


License
-------

Hivy is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).
