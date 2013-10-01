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


import yaml
import os
import glob

import unittest
import utils


class TestUtils(unittest.TestCase):
    '''
    Utilities used by pencil client
    '''
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_load_yaml(self):
        content = utils.load_yaml("client/sample-hivy.yml")
        assert isinstance(content, dict)
        assert "cells" in content
        assert "project" in content

    def test_load_yaml_no_yaml(self):
        content = utils.load_yaml("client/broken-sample-hivy.yml")
        assert "error" in content
        assert isinstance(content["error"], yaml.scanner.ScannerError)

    def test_load_yaml_no_file(self):
        content = utils.load_yaml("fake-sample-hivy.yml")
        assert "error" in content
        assert isinstance(content["error"], IOError)

    def test_ui(self):
        utils.fail("TestUI")
        utils.success("TestUI")
        utils.log("TestUI")

    def test_store_certificate(self):
        utils.store_certificate("- BEGIN CERTIFICATE -\nfakecertificate", ".")
        cacrt = filter(os.path.isfile, glob.glob('ca-*.crt'))
        assert len(cacrt)
        os.remove(cacrt[0])

    def test_store_certificate_wrong_path(self):
        try:
            utils.store_certificate("-- BEGIN CERTIFICATE --\nfakecertificate",
                                    "/fake/path")
        except ValueError, e:
            print("This is normal behaviour: ", e)


if __name__ == '__main__':
    unittest.main()
