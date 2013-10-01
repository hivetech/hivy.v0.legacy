Hivy
====

Setup
-----

First make sure [etcd](https://github.com/coreos/etcd) is installed and available in your $PATH.

```bash
$ git clone https://github.com/hivetech/hivy.go
$ make
$ ./hivy --help
```
TODO: Automatic user creation


Usage
-----

```bash
$ ./hivy -d node -n master --verbose --cpuprofile profile -f
$ curl --user xav:boss http://127.0.0.1:8080/dummy  # Test your installation
$ # In another terminal
$ ./client/pencil login
$ ./client/pencil configure --app quantrade --config client/sample-hivy.yml
```


API
---

```
GET /dummy/
GET /login/{user}
GET /deploy/{project}
```


Configuration storage topology
------------------------------

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
