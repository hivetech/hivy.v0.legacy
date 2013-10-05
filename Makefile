# Makefile for juju-core.
# vim:ft=make

PROJECT="github.com/hivetech/hivy"
CHARMSTORE?="${HOME}/charms"

all: install extras-dev

tests: check coverage 
	@echo "Done."

coverage:
	gocov test github.com/hivetech/hivy github.com/hivetech/hivy/filters github.com/hivetech/hivy/endpoints github.com/hivetech/hivy/security | gocov report

check:
	go build
	go test -short -gocheck.v

# Install packages required to develop Juju and run tests.
local-install:
	go get -u github.com/mattn/gom
	#FIXME Include as well hivy sub-packages
	test -f Gomfile || gom gen gomfile
	gom install
	gom build

install:
	git clone https://github.com/coreos/etcd.git /tmp/etcd
	cd /tmp/etcd/ && ./build
	test -f /tmp/etcd/etcd && cp /tmp/etcd/etcd ${GOPATH}/bin
	cd -

	cat Gomfile | sed -e s/gom\ // | xargs go get -u
	go install

run:
	go build
	./hivy -d node -n master --verbose

extras-dev:
	sudo apt-get install python-pip
	sudo pip install -U httpie

	go get -u github.com/mattn/gom
	go get -u github.com/gophertown/looper
	go get -u launchpad.net/gocheck
	go get -u github.com/remogatto/prettytest
	go get -u github.com/axw/gocov/gocov

watch:
	pgrep --count etcd > /dev/null || etcd -n master -d node -v &
	looper -debug
	killall etcd

init:
	pgrep --count etcd > /dev/null || etcd -n master -d node -v &
	curl -L http://127.0.0.1:4001/v1/keys/hivy/charmstore -d value="${CHARMSTORE}"
	curl -L http://127.0.0.1:4001/v1/keys/hivy/security/admin/password -d value="root"
	curl -L http://127.0.0.1:4001/v1/keys/hivy/security/admin/methods/GET/createuser -d value="1"
	killall etcd

doc:
	godoc -http=:6060

.PHONY: install format check
