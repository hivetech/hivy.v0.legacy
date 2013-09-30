#!/bin/bash

TEST_USER="xav"
TEST_PASSWORD="boss"
TEST_PROJECT="quantrade" 

#configs setter
curl -L --user $TEST_USER:$TEST_PASSWORD http://127.0.0.1:8080/dummy
