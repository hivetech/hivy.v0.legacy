#!/bin/bash

juju destroy-service patate-coin-hivelab
juju destroy-service user-project-hivelab
juju destroy-service user-project-mysql
juju destroy-service user-project-wordpress

juju status
