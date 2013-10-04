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


import sys
import yaml
import os
import clint.textui as textui
import inspect


class PencilError(Exception):
    def __init__(self, **kwargs):
        super(PencilError, self).__init__(kwargs)
        self.code = kwargs.get('errorCode')
        self.message = kwargs.get('message')
        self.cause = kwargs.get('cause')


def fail(message):
    textui.puts(textui.colored.red("[{}] ** Error: {}".format(
        inspect.stack()[1][3], message)))


def die(message):
    sys.exit(textui.colored.red("[{}] ** Error: {}".format(
        inspect.stack()[1][3], message)))


def success(message):
    textui.puts(textui.colored.green(message))


def log(message):
    textui.puts(textui.colored.blue(message))


def load_yaml(filepath):
    try:
        return yaml.load(open(filepath, "r"))
    except yaml.scanner.ScannerError, e:
        fail(e)
        return {'error': e}
    except IOError, e:
        fail(e)
        return {'error': e}
    return {}


def store_certificate(certificate, path="/etc/ssl/certs"):
    '''
    The provided certificate will be used later to allow etcd transactions.
    So we need to store it somewhere we can remember.
    '''
    cert_name = "ca-{}.crt".format(hash('hivy'))
    if not os.path.isdir(path):
        raise ValueError("Path must exist")
    with open(os.path.join(path, cert_name), 'w') as fd:
        fd.write(certificate)
