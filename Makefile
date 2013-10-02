# Makefile for juju-core.
# vim:ft=make

PROJECT="github.com/hivetech/hivy"

all: install extras-dev test

tests: check-server check-client
	@echo "Done."

# Run tests.
check-server:
	ps -ef | grep etcd | grep -v etcd || etcd -n master -d node -v &
	go test -test.v
	killall etcd

check-client:
	ps -ef | grep hivy | grep -v hivy || ./hivy -n master -d node --verbose &
	nosetests --verbose --with-progressive client
	killall hivy

# Install packages required to develop Juju and run tests.
local-install:
	go get -u github.com/mattn/gom
	#FIXME Include as well hivy sub-packages
	test -f Gomfile || gom gen gomfile
	gom install
	gom build

install:
	apt-get install python-pip
	pip install -r client/requirements.txt
	cat Gomfile | sed -e s/gom\ // | xargs go get -u
	go install

run:
	go build
	./hivy -d node -n master --verbose

extras-dev:
	go get -u github.com/mattn/gom
	go get -u github.com/gophertown/looper
	go get -u launchpad.net/gocheck
	go get -u github.com/remogatto/prettytest

watch:
	ps -ef | grep etcd | grep -v etcd || etcd -n master -d node -v &
	looper -debug
	killall etcd

doc:
	godoc -http=:6060

.PHONY: install format check
