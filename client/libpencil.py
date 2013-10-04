#! /usr/bin/env python
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#
# Copyright 2013 Xavier Bruhiere
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


import etcd
import yaml
import os
import requests
import clint.textui as textui
from random import random
import time

import utils


class Pencil(etcd.Etcd):
    '''
    Main client class holding methods that will talk to the Hive server.
    '''
    __version__ = "0.0.1"
    __credentials__ = os.path.expanduser("~/.hivy.yml")
    user = None
    password = None

    def __init__(self, **kwargs):
        self.project = kwargs.pop("project", os.getcwd().split('/')[-1])
        #FIXME kwargs['port'] is for etcd, server port is hard-coded,
        #      ip are common
        self.host = kwargs.get("host", '127.0.0.1')
        self.port = 8080

        if os.path.exists(self.__credentials__):
            credentials = utils.load_yaml(self.__credentials__)
            self.user = credentials["user"]
            self.password = credentials["password"]

        try:
            etcd.Etcd.__init__(self, **kwargs)
        except requests.exceptions.ConnectionError, e:
            utils.die(e)

    def projectpath(self):
        return os.path.join(self.user, self.project)

    def charmpath(self, name):
        return os.path.join(self.projectpath(), name)

    def _hivy_request(self, path, data={}):
        try:
            result = requests.get('http://{}:{}/{}'.format(
                                  self.host, self.port, path),
                                  params=data,
                                  auth=(self.user, self.password))
        except requests.exceptions.ConnectionError, e:
            utils.die(e)
        try:
            return result.json()
        except:
            utils.fail(result.content)
        return {'error': result.content}

    def configure(self, namespace, key, value):
        #TODO if isinstance(value, dict):
        if isinstance(value, list):
            value = ",".join(value)

        key_etcd_path = os.path.join(namespace, key)

        try:
            feedback = self.set(key_etcd_path, value)
        except etcd.EtcdError, e:
            raise(e)
        return feedback

    def _store_credentials(self, data):
        with open(self.__credentials__, 'w') as fd:
            fd.write(yaml.dump(data, default_flow_style=False))

    def up(self):
        '''
        Deploy configured cell
        '''
        return self._hivy_request('/juju/deploy/{}'.format(self.project),
                                  {'user': self.user})

    def login(self, username=None, password=None):
        '''
        Use https internal authentification mechanism to send username and
        password to the server, and fetch back a certificate.
        '''
        is_ok = False
        utils.log("Please submit your Hive credentials")
        if not username:
            username = raw_input("\tUsername  ")
        if not password:
            password = raw_input("\tPassword  ")

        for _ in textui.progress.bar(range(100)):
            time.sleep(random() * 0.02)
        try:
            result = requests.get('http://{}:{}/login/'.
                                  format(self.host, self.port),
                                  auth=(username, password))
        except requests.exceptions.ConnectionError, e:
            utils.die(e)

        certificate = result.text
        if certificate.find("CERTIFICATE") > 0:
            utils.success("Successfully logged in.")
            self._store_credentials({'user': username, 'password': password})

            utils.log("Certificate provided, storing it.")
            utils.store_certificate(certificate, path=".")
            is_ok = True
        else:
            utils.fail("Login failed: no certificate returned. ({})".format(
                certificate))
        return is_ok

    def help(self, api=""):
        #TODO Pretty printing
        print(self._hivy_request("/help/{}".format(api)))
