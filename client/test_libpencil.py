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


import os
import unittest
import etcd
from libpencil import Pencil


class TestUtils(unittest.TestCase):
    '''
    Pencil tests
    '''
    def setUp(self):
        self.default_project = "quantrade"
        self.default_host = "127.0.0.1"
        self.default_port = 4001
        self.default_user = "xav"
        self.default_password = "boss"

        self.c = Pencil(project=self.default_project,
                        host=self.default_host,
                        port=self.default_port,
                        follow_leader=True, autostart=True)

    def tearDown(self):
        pass

    def test_pencil_init(self):
        assert(Pencil(project=self.default_project,
                      host=self.default_host,
                      port=self.default_port,
                      follow_leader=True, autostart=True))

    def test_projectpath(self):
        path = self.c.projectpath()
        self.assertEqual(path, "{}/{}".format(
            self.default_user, self.default_project))

    def test_charmpath(self):
        path = self.c.charmpath("test")
        self.assertEqual(path, "{}/{}/{}".format(
            self.default_user, self.default_project, "test"))

    def test_configure(self):
        #TODO Run etcd
        path = '{}/{}'.format(self.default_user, self.default_project)
        feedback = self.c.configure(path, "session", "test")
        assert isinstance(feedback, etcd.etcd.EtcdSet)

    def test_up(self):
        feedback = self.c.up()
        #FIXME Should be json
        assert isinstance(feedback, unicode)

    def test__store_credentials(self):
        self.c.__credentials__ = "test-credentials.yml"
        self.c._store_credentials("foo: bar")
        os.remove("test-credentials.yml")

    def test__store_credentials_not_yaml(self):
        self.c.__credentials__ = "test-credentials.yml"
        self.c._store_credentials("foo")
        os.remove("test-credentials.yml")

    def test_login_with_credentials(self):
        is_ok = self.c.login(self.default_user, self.default_password)
        assert(is_ok)

    def test_login_with_wrong_credentials(self):
        is_ok = self.c.login(self.default_user, "fake")
        assert(not is_ok)


if __name__ == '__main__':
    unittest.main()
