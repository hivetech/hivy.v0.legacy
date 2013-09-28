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
import sys
import yaml
import os
import requests
import clint.textui as textui
from random import random
import time
import inspect


def fail(message):
    textui.colored.red("[{}] ** Error: {}".format(
        inspect.stack()[1][3], message))


def die(message):
    sys.exit(textui.colored.red("[{}] ** Error: {}".format(
        inspect.stack()[1][3], message)))


def success(message):
    textui.puts(textui.colored.green(message))


def log(message):
    textui.puts(textui.colored.blue(message))


def load_yaml(filepath):
    return yaml.load(open(filepath, "r"))


def store_certificate(certificate, path="/etc/ssl/certs"):
    '''
    The provided certificate will be used later to allow etcd transactions.
    So we need to store it somewhere we can remember.
    '''
    cert_name = "ca{}.crt".format(hash('hivy'))
    with open(os.path.join(path, cert_name), 'w') as fd:
        fd.write(certificate)


class Pencil(etcd.Etcd):
    '''
    Main client class holding methods that will talk to the Hive server.
    '''
    __version__ = "0.0.1"
    __credentials__ = os.path.expanduser("~/.hivy.yml")
    user = None
    password = None

    def __init__(self, **kwargs):
        try:
            etcd.Etcd.__init__(self)
        except requests.exceptions.ConnectionError, e:
            die(e)

        self.project = kwargs.get("project", os.getcwd().split('/')[-1])

        if os.path.exists(self.__credentials__):
            credentials = load_yaml(self.__credentials__)
            self.user = credentials["user"]
            self.password = credentials["password"]

    def projectpath(self):
        return os.path.join(self.user, self.project)

    def charmpath(self, name):
        return os.path.join(self.projectpath(), name)

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

    def up(self):
        '''
        Deploy configured cell
        '''
        try:
            #FIXME Hard coded
            result = requests.get('http://127.0.0.1:8080/deploy',
                                  params={'user': self.user,
                                          'project': self.project},
                                  auth=(self.user, self.password))
        except requests.exceptions.ConnectionError, e:
            die(e)
        #FIXME If not json output, crash
        return result.json()

    def _store_credentials(self, data):
        with open(self.__credentials__, 'w') as fd:
            fd.write(yaml.dump(data, default_flow_style=False))

    def login(self, extra="test"):
        '''
        Use https internal authentification mechanism to send username and
        password to the server, and fetch back a certificate.
        '''
        is_ok = False
        log("Please submit your Hive credentials")
        username = raw_input("\tUsername  ")
        password = raw_input("\tPassword  ")

        self._store_credentials({'user': username, 'password': password})

        for _ in textui.progress.bar(range(100)):
            time.sleep(random() * 0.02)
        try:
            #FIXME Hard coded
            result = requests.get('http://127.0.0.1:8080/login/{}'.
                                  format(extra), auth=(username, password))
        except requests.exceptions.ConnectionError, e:
            die(e)

        success("Successfully logged in.")
        if 'Cacrt' in result.json():
            log("Certificate provided, storing it.")
            store_certificate(result.json()['Cacrt'], path=".")
            is_ok = True
        else:
            fail("Login failed: no certificate returned.")
        return is_ok


#TODO save user/project input as default
def main(args):
    '''
    Main entry for final user
    Use arguments parsed by docopt
    '''

    c = Pencil(project=args['--app'])

    if args['configure']:
        if args['--config']:
            configuration = load_yaml(args['--config'])

            for k, v in configuration.items():
                if k == 'cells':
                    services = []
                    for cell in configuration['cells']:
                        services.append(cell['charm'])
                        for k_cell, v_cell in cell.items():
                            print(c.configure(
                                c.charmpath(cell['charm']), k_cell, v_cell))
                    print(c.configure(
                        c.projectpath(), 'services', ','.join(services)))
                else:
                    print(c.configure(c.projectpath(), k, v))

        if args['--key']:
            keyvalue = args['--key'].split(':')
            print(c.configure(c.projectpath(), keyvalue[0], keyvalue[1]))

    elif args['up']:
        if args['--config'] or args['--key']:
            raise(NotImplementedError)
        print(c.up())

    elif args['login']:
        #TODO Use user and/or password if provided
        print(c.login())

    elif args['inspect']:
        raise(NotImplementedError)
