# Makefile for juju-core.
# vim:ft=make

PROJECT="github.com/hivetech/hivy"

all: install extras-dev test

tests: check-server check-client coverage 
	@echo "Done."

coverage:
	#FIXME Does not prevent etcd to run
	ps -ef | grep etcd | grep -v etcd || etcd -n master -d node -v &
	#FIXME gocov test github.com/hivetech/hivy | gocov report
	gocov test github.com/hivetech/hivy/security | gocov report
	killall etcd

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
	#mkdir build
	#git clone https://github.com/coreos/etcd.git build/etcd
	#./build/etcd/build
	#test -f ./build/etcd/etcd && cp ./build/etcd/etcd ${GOPATH}/bin

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
	go get -u github.com/axw/gocov/gocov

watch:
	ps -ef | grep etcd | grep -v etcd || etcd -n master -d node -v &
	looper -debug
	killall etcd

doc:
	godoc -http=:6060

.PHONY: install format check
