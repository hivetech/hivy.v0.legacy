#! /usr/bin/env python
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#
# Copyright © 2013 xavier <xavier@laptop-300E5A>
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
import sys
import yaml
import os
import requests
import clint.textui as textui
from random import random
import time


class Pencil(etcd.Etcd):
    '''
    Main client class holding methods that will talk to the Hive server.
    '''
    __version__ = "0.0.1"

    def __init__(self, user, project):
        try:
            etcd.Etcd.__init__(self)
        except requests.exceptions.ConnectionError, e:
            sys.exit(
                textui.colored.red("[Pencil.__init__] ** Error: {}".format(e)))
        self.user = user
        self.project = project
        # Used to store arrays of value at a given key
        self.cell_index = {}

    def SetCellConfig(self, cell, key, values):
        '''
        Store cell specific configuration for a given project
        '''
        # Feedback stores etcd response per request
        feedbacks = {}
        # Process single string value and array of values
        # So we sanitize it this way:
        if isinstance(values, str):
            values = [values]

        # Record "key: value" to etcd one by one
        #set USER/PROJECT/config/<key>/<index> -d value=<value>
        for value in values:
            # Firstly manage index
            if cell in self.cell_index:
                if key in self.cell_index[cell]:
                    self.cell_index[cell] += 1
                else:
                    self.cell_index[cell][key] = 0
            else:
                self.cell_index[cell][key] = 0

            key_etcd_path = os.path.join(self.user, self.project,
                                         cell, "config", key, str(value))
            try:
                feedbacks[value] = self.set(key_etcd_path, value)
            except etcd.EtcdError, e:
                raise(e)
        return feedbacks

    def SetProjectConfig(self, key, value):
        key_etcd_path = os.path.join(self.user, self.project, key)
        try:
            feedback = self.set(key_etcd_path, value)
        except etcd.EtcdError, e:
            raise(e)
        return feedback

    def Up(self):
        '''
        Deploy configured cell
        '''
        return self.SetProjectConfig("state", "run")

    def read_config(self, filepath):
        return yaml.load(open(filepath, "r"))

    def _store_certificate(self, certificate, path="/etc/ssl/certs"):
        '''
        The provided certificate will be used later to allow etcd transactions.
        So we need to store it somewhere we can remember.
        '''
        cert_name = "ca{}.crt".format(hash(self.user + self.project))
        with open(os.path.join(path, cert_name), 'w') as fd:
            fd.write(certificate)

    def login(self, extra="test"):
        '''
        Use https internal authentification mechanism to send username and
        password to the server, and fetch back a certificate.
        '''
        is_ok = False
        textui.puts(textui.colored.blue("Please submit your Hive credentials"))
        username = raw_input("\tUsername  ")
        password = raw_input("\tPassword  ")
        for _ in textui.progress.bar(range(100)):
            time.sleep(random() * 0.02)
        try:
            #FIXME Hard coded
            result = requests.get('http://127.0.0.1:8080/login/{}'.
                                  format(extra), auth=(username, password))
        except requests.exceptions.ConnectionError, e:
            sys.exit(
                textui.puts(textui.colored.red(
                    "[Pencil.login] ** Error: {}".format(e))))

        textui.puts(textui.colored.green("Successfully logged in."))
        if 'Cacrt' in result.json():
            textui.puts(textui.colored.blue(
                "Certificate provided, storing it."))
            self._store_certificate(result.json()['Cacrt'], path=".")
            is_ok = True
        else:
            textui.puts(textui.colored.red(
                "Login failed: no certificate returned."))
        return is_ok


def debug_etcd(feedback):
        print "Index:", feedback.index
        print "Newkey:", feedback.newKey
        print "Previous value:", feedback.prevValue
        print "Expiration:", feedback.expiration
        print

#TODO save user/project input as default
if __name__ == '__main__':
    config_file = "sample-hivy.yml"
    user = "xavier"
    project = "hivetech"

    c = Pencil(user, project, config=config_file)

    if sys.argv[1] == "setconfig":
        configuration = c.read_config(config_file)
        for k, v in configuration.items():
            debug_etcd(c.SetProjectConfig(k, v))

        print
        debug_etcd(c.SetProjectConfig("state", "loaded"))

    elif sys.argv[1] == "up":
        debug_etcd(c.Up())

    elif sys.argv[1] == "login":
        print(c.login())
