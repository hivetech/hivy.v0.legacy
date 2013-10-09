#!/bin/bash

# Copyright 2013 tsuru authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

status=0
out=`gofmt -s -l .`
if [ "${out}" != "" ]
then
    echo "ERROR: there are files that need to be formatted with gofmt"
    echo
    echo "Files:"
    for file in $out
    do
        echo "- ${file}"
        #TODO With command line: gofmt -w ${file}
    done
    status=1
else
    echo "Code format is clean"
fi

`go vet ./... > .vet 2>&1`
out=`cat .vet`
if [ "${out}" != "" ]
then
    echo
    echo "ERROR: go vet failures:"
    echo
    cat <<END
${out}
END
    status=1
fi

rm .vet || true
exit $status

