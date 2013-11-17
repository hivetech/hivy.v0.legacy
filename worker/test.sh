#! /bin/bash
#
# test.sh
# Copyright (C) 2013 xavier <xavier@laptop-300E5A>
#
# Distributed under terms of the MIT license.
#

redis-cli -r 100 RPUSH resque:queue:forker '{"class":"Hivy","args":["ls", "-al"]}'
