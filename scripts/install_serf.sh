#!/bin/sh
#
# This script installs and configures the Serf agent that runs on
# every node. As with the other scripts, this should probably be done with
# formal configuration management, but a shell script is simple as well.

set -e

version=$1

sudo apt-get install -y unzip

# Download and install Serf
cd /tmp
until wget -O serf.zip https://dl.bintray.com/mitchellh/serf/${version}.zip; do
    sleep 1
done
unzip serf.zip
sudo mv serf /usr/local/bin/serf
