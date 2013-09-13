#!/bin/bash

testUSER="patate"
testPROJECT="coin" 

testUSER2="user"
testPROJECT2="project"

#configs setter
curl -L http://127.0.0.1:4001/v1/keys/$testUSER/$testPROJECT/config/series/0 -d value="precise"
curl -L http://127.0.0.1:4001/v1/keys/$testUSER/$testPROJECT/config/cell/0 -d value="hivelab"
curl -L http://127.0.0.1:4001/v1/keys/$testUSER/$testPROJECT/config/version/0 -d value="latest"
curl -L http://127.0.0.1:4001/v1/keys/$testUSER/$testPROJECT/config/services/0 -d value="none"
curl -L http://127.0.0.1:4001/v1/keys/$testUSER/$testPROJECT/config/packages/0 -d value="none"

curl -L http://127.0.0.1:4001/v1/keys/$testUSER2/$testPROJECT2/config/series/0 -d value="precise"
curl -L http://127.0.0.1:4001/v1/keys/$testUSER2/$testPROJECT2/config/cell/0 -d value="hivelab"
curl -L http://127.0.0.1:4001/v1/keys/$testUSER2/$testPROJECT2/config/version/0 -d value="latest"
curl -L http://127.0.0.1:4001/v1/keys/$testUSER2/$testPROJECT2/config/services/0 -d value="mysql"
curl -L http://127.0.0.1:4001/v1/keys/$testUSER2/$testPROJECT2/config/services/1 -d value="wordpress"
curl -L http://127.0.0.1:4001/v1/keys/$testUSER2/$testPROJECT2/config/packages/0 -d value="vim"
curl -L http://127.0.0.1:4001/v1/keys/$testUSER2/$testPROJECT2/config/packages/1 -d value="zsh"

#config done
#NB the state "config loaded" is only used as pre-state value of any other run-like cmd
curl -L http://127.0.0.1:4001/v1/keys/$testUSER/$testPROJECT/state -d value="loaded"

#run
curl -L http://127.0.0.1:4001/v1/keys/$testUSER/$testPROJECT/state -d value="run"

#multi user/statement test
curl -L http://127.0.0.1:4001/v1/keys/$testUSER2/$testPROJECT2/state -d value="loaded"
curl -L http://127.0.0.1:4001/v1/keys/$testUSER2/$testPROJECT2/state -d value="run"
