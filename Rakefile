# = Hivy App Rake - Unide router
#
# Author::    Xavier Bruhiere
# Copyright:: (c) 2013, Xavier Bruhiere
# License::   Apache 2.0
#
# Feedback appreciated: xavier.bruhiere@gmail.com

# https://github.com/michaeldv/awesome_print
require "awesome_print"

logs = "/tmp/hivy.rake"

task :default => [:install]

desc "Install hivy application"
task :install => [ "install:deps",  "install:etcd"]

namespace :install do
    desc "Install app dependencies"
    task :deps do
      msg "generate go dependencies"
      sh "test -f Gomfile || gom gen gomfile"
      msg "for dev time, remove hivetech libs from deps"
	    sh "sed -i '/hivetech/d' Gomfile"
      msg "install go dependencies"
	    sh "cat Gomfile | sed -e s/gom\\ // | xargs go get -u"
      msg "install boxcars proxy"
      sh "go get -u github.com/hivetech/boxcars/boxcars"

      msg "compile and install workers"
      sh "cd worker && go install"
      msg "compile and install hivy app"
      cmd = "go install"
      result = system(cmd)
      raise("optparse installation failed..  msg: #{$?}") unless result
    end

    desc "Install configuration storage, etcd"
    task :etcd do
      msg "fetch etcd code"
	    sh "test -d /tmp/etcd || git clone https://github.com/coreos/etcd.git /tmp/etcd"
      msg "build and install it"
	    sh "cd /tmp/etcd/ && ./build"
	    sh "test -f /tmp/etcd/etcd && cp /tmp/etcd/etcd ${GOPATH}/bin"
    end
end

namespace :app do
    desc "update and run the app"
    task :run do
      msg "install workers"
      sh "cd worker && go install"
      msg "re-compile and install"
      sh "go install"
      msg "run application"
      sh "forego start"
    end
end

private

def msg(text)
    ap "  => rake: #{text}"
end
