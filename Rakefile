# = Hivy App Rake - Unide router
#
# Author::    Xavier Bruhiere
# Copyright:: (c) 2013, Xavier Bruhiere
# License::   Apache 2.0
#
# Feedback appreciated: xavier.bruhiere@gmail.com

# https://github.com/michaeldv/awesome_print
#require "awesome_print"

logs = "/tmp/hivy.rake"
hivy_charm_repo = "https://github.com/hivetech/cells.git"
charmstore = ENV["HOME"] + "/cells"
goworkspace = ENV["HOME"] + "/goworkspace"

task :default => [:install]

desc "Install hivy application"
task :install => [ "install:deps", "install:vendor", "install:compile" ]

desc "update hivy dependencies"
task :update => [ "install" ]

desc "update hivy dependencies"
task :tests => [ "tests:check", "tests:coverage", "tests:style" ]

namespace :install do
  desc "Install go version 1.1.1"
  task :go do
    msg "Install go version 1.1.1"
    sh "sudo apt-get install -y python-software-properties"
    sh "sudo add-apt-repository -y -s ppa:duh/golang"
    sh "sudo apt-get update"
    sh "sudo apt-get install -y golang"
    sh "go version"
    msg "Create workspace"
    sh "mkdir #{goworkspace}"
  end

  desc "Install app dependencies"
  task :deps do
    msg "Install package dependencies"
    sh "sudo apt-get install bzr mercurial"
    msg "install go dependencies"
    sh "go get -v -u github.com/benmanns/goworker"
    sh "go get -v -u github.com/ddollar/forego"
    sh "go get -v -u github.com/codegangsta/cli"     
    msg "install boxcars proxy"
    #FIXME Need stable version: sh "go get -u github.com/hivetech/boxcars/boxcars"
    sh "go get -v -u github.com/azer/boxcars/boxcars"     
    msg "installs deps manager: gom"
    sh "go get -u github.com/mattn/gom"
    msg "generate go dependencies"
    sh "test -f Gomfile || gom gen gomfile"
    msg "for dev time, remove hivetech libs from deps"
    sh "sed -i '/hivetech/d' Gomfile"
    sh "cat Gomfile | sed -e s/gom\\ // | xargs go get -v -u"
  end

  desc "Compile workers and hivy"
  task :compile do
    msg "Compile workers and hivy"
    msg "compile and install workers"
    sh "cd worker && go install"
    msg "compile and install hivy app"
    cmd = "cd hivy && go install"
    result = system(cmd)
    raise("optparse installation failed..  msg: #{$?}") unless result
  end

  desc "Install hivy backends dependencies"
  task :vendor do
    msg "fetch etcd code"
    sh "test -d /tmp/etcd || git clone https://github.com/coreos/etcd.git /tmp/etcd"
    msg "build and install it"
    sh "cd /tmp/etcd/ && ./build"
    sh "test -f /tmp/etcd/etcd && cp /tmp/etcd/etcd ${GOPATH}/bin"

    msg "install serf, service orchestration and monitoring"
    sh "./scripts/install_serf.sh 0.2.1_linux_amd64"

    msg "install redis, jobs queue"
    sh "sudo apt-get install -y redis-server"

    msg "install lxc containers technology"
    sh "sudo apt-get install -y lxc"

    msg "Download hivy charms"
    sh "[ -d #{charmstore} ] || git clone #{hivy_charm_repo} #{charmstore}"

    #FIXME go get ends up with an error
    msg "download juju-core"
    sh "go get -v launchpad.net/juju-core/..."
  end

  desc "requested packages for hivy app and scripts"
  task :client do
    msg "Installing extra client dependencies"
    sh "sudo pip install -U httpie"

    #TODO needs npm and node to be installed
    sh "npm install -g underscore-cli"
  end

  desc "Uses gom to install hivy localy"
  task :local do
    sh "go get -u github.com/mattn/gom"
    sh "test -f Gomfile || gom gen gomfile"
    sh "gom install"
    sh "gom build"
  end
end

namespace :app do
  desc "setup minimum admin user to manage hivy"
  task :init do
	  sh "pgrep --count etcd > /dev/null || etcd -n master -d hivy/.conf -v &"
	  sh "sleep 5"
    sh "curl -L http://127.0.0.1:4001/v1/keys/hivy/charmstore -d value=#{charmstore}"
    sh "curl -L http://127.0.0.1:4001/v1/keys/hivy/security/admin/password -d value=root"
    sh "curl -L http://127.0.0.1:4001/v1/keys/hivy/security/admin/methods/PUT/v0/methods/user -d value=1"
    sh "curl -L http://127.0.0.1:4001/v1/keys/hivy/security/admin/methods/GET/v0/methods/help -d value=1"
    sh "curl -L http://127.0.0.1:4001/v1/keys/hivy/security/admin/methods/GET/v0/methods/dummy -d value=1"
    sh "curl -L http://127.0.0.1:4001/v1/keys/hivy/mapping/port -d value=49153"
    sh "killall etcd"
  end

  desc "updates and run the app"
  task :run do
    msg "install workers"
    sh "cd worker && go install"
    msg "re-compile and install"
    sh "cd hivy && go install"
    msg "run application"
    sh "cd hivy && forego start"
  end

  desc "start documentation server on localhost:6060"
  task :doc do
    msg "connect at localhost:6060 to see package documentation"
    sh "godoc -http=:6060"
  end
end

namespace :utils do
  desc "removes temporary or useless files"
  task :clean do
    sh "rm -rf *.test hivy/.profile"
  end

  desc "sums up whole application lines"
  task :count do
	  sh "find . -name \"*.go\" -print0 | xargs -0 wc -l"
  end
end

namespace :tests do
  desc "installs tests dependencies"
  task :deps do
    # I'm testing pretty much everything about go...
    sh "go get -u github.com/gophertown/looper"
    sh "go get -u launchpad.net/gocheck"
    sh "go get -u github.com/remogatto/prettytest"
    sh "go get -u github.com/axw/gocov/gocov"
    sh "go get -u github.com/golang/lint/golint"
    sh "go get -u github.com/davecheney/profile"
    sh "go get -u github.com/mreiferson/go-httpclient"
  end

  desc "runs tests"
  task :check do
    msg "compile libs"
    sh "go test -i"
    msg "fire up tests"
    sh "go test -short -gocheck.v"
  end

  desc "computes project coverage"
  task :coverage do
    msg "computes project coverage"
    sh "gocov test github.com/hivetech/hivy github.com/hivetech/hivy/security | gocov report"
  end

  desc "analyses coding style"
  task :style do
	  msg "style.govet"
	  sh "find . -name \"*.go\" | xargs go tool vet -all -v"
	  msg "style.golint"
	  sh "find . -name \"*.go\" | xargs golint -min_confidence=0.8"
  end

  desc "continuous testing"
  task :watch do
    msg "monitoring files for continuous testing"
    sh "looper -debug"
  end
end

private

def msg(text)
    #ap "  => rake: #{text}"
    puts "  => rake: #{text}"
end
