# Makefile for juju-core.
# vim:ft=make

PROJECT="github.com/hivetech/hivy"
CHARMSTORE?="${HOME}/charms"

all: install extras-dev

update: install extras-dev

tests: check coverage style
	@echo "Done."

coverage:
	@echo "\t[ make ] ===========>  Coverage"
	gocov test github.com/hivetech/hivy github.com/hivetech/hivy/endpoints github.com/hivetech/hivy/security | gocov report

check:
	go build
	go test -i
	@echo -e "\t[ make ] ===========>  Tests"
	go test -short -gocheck.v

style:
	@echo "\t[ make ] ===========>  Style.govet"
	find . -name "*.go" | xargs go tool vet -all -v
	@echo "\t[ make ] ===========>  Style.golint"
	find . -name "*.go" | xargs golint -min_confidence=0.8 

# Install packages required to develop Juju and run tests.
local-install:
	go get -u github.com/mattn/gom
	test -f Gomfile || gom gen gomfile
	sed -i '/hivetech/d' Gomfile
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
	go install 
	forego start

extras-dev:
	sudo apt-get install python-pip
	sudo pip install -U httpie

	npm install -g underscore-cli

	# I'm testing pretty much everything about go...
	go get -u github.com/mattn/gom
	go get -u github.com/gophertown/looper
	go get -u launchpad.net/gocheck
	go get -u github.com/remogatto/prettytest
	go get -u github.com/axw/gocov/gocov
	go get -u github.com/golang/lint/golint
	go get -u github.com/davecheney/profile
	go get -u github.com/ddollar/forego

watch:
	looper -debug

init:
	pgrep --count etcd > /dev/null || etcd -n master -d node -v &
	curl -L http://127.0.0.1:4001/v1/keys/hivy/charmstore -d value="${CHARMSTORE}"
	curl -L http://127.0.0.1:4001/v1/keys/hivy/security/admin/password -d value="root"
	curl -L http://127.0.0.1:4001/v1/keys/hivy/security/admin/methods/GET/createuser -d value="1"
	killall etcd

doc:
	godoc -http=:6060

clean:
	rm -rf profile/* *.test hivy

.PHONY: install format check
