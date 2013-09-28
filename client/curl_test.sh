#!/bin/bash

TEST_USER="patate"
TEST_PROJECT="coin" 

#configs setter
curl -L http://127.0.0.1:4001/v1/keys/$TEST_USER/$TEST_PROJECT/config/series/0 -d value="precise"
curl -L http://127.0.0.1:4001/v1/keys/$TEST_USER/$TEST_PROJECT/config/cell/0 -d value="hivelab"
curl -L http://127.0.0.1:4001/v1/keys/$TEST_USER/$TEST_PROJECT/config/version/0 -d value="latest"
curl -L http://127.0.0.1:4001/v1/keys/$TEST_USER/$TEST_PROJECT/config/services/0 -d value="none"
curl -L http://127.0.0.1:4001/v1/keys/$TEST_USER/$TEST_PROJECT/config/packages/0 -d value="none"

#config done
#NB the state "config loaded" is only used as pre-state value of any other run-like cmd
curl -L http://127.0.0.1:4001/v1/keys/$TEST_USER/$TEST_PROJECT/state -d value="loaded"

#run
curl -L http://127.0.0.1:4001/v1/keys/$TEST_USER/$TEST_PROJECT/state -d value="run"

#multi user/statement test
curl -L http://127.0.0.1:4001/v1/keys/$TEST_USER2/$TEST_PROJECT2/state -d value="loaded"
curl -L http://127.0.0.1:4001/v1/keys/$TEST_USER2/$TEST_PROJECT2/state -d value="run"
