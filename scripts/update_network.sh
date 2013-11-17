#! /bin/bash
#
# setup_iptables.sh
# Copyright (C) 2013 xavier <xavier@laptop-300E5A>
#
# Distributed under terms of the MIT license.

# Inspiration : https://github.com/dotcloud/docker/blob/master/iptables/iptables.go

# const (
  export RULES_NAME="JUJU"
  export BRIDGE="lxcbr0"
  export HOST_IP=$(hostname -I | awk '{print $1}')
  export PROTOCOL="TCP"
  export portRangeStart=49153
  export portRangeEnd=65535
# )

#TODO Mobile notification

function setup_new_chain() {
  sudo iptables -t nat -N ${RULES_NAME}
  sudo iptables -t nat -A PREROUTING -m addrtype --dst-type LOCAL -j ${RULES_NAME}
  sudo iptables -t nat -A OUTPUT -m addrtype --dst-type LOCAL ! --dst 127.0.0.0/8 -j ${RULES_NAME}

  etcdctl set hivy/mapping/port ${portRangeStart}
}

function clean_iptables() {
  sudo iptables -t nat -D PREROUTING -m addrtype --dst-type LOCAL -j ${RULES_NAME}
  sudo iptables -t nat -D OUTPUT -m addrtype --dst-type LOCAL --dst 127.0.0.0/8 -j ${RULES_NAME}
  sudo iptables -t nat -D OUTPUT -m addrtype --dst-type LOCAL -j ${RULES_NAME}

  sudo iptables -t nat -D PREROUTING -j ${RULES_NAME}
  sudo iptables -t nat -D OUTPUT -j ${RULES_NAME}

  sudo iptables -t nat -F ${RULES_NAME}
  sudo iptables -t nat -X ${RULES_NAME}
}

function map_virtual_port() {
  sudo iptables -t nat -A ${RULES_NAME} -p ${PROTOCOL} -d ${HOST_IP} --dport ${HOST_PORT} ! \
    -i ${BRIDGE} -j DNAT --to-destination ${VIRTUAL_IP}:22 
}

function allocate_port() {
  port=$(etcdctl get hivy/mapping/port)
  etcdctl set hivy/mapping/port $(($port+1))
  if [ ${port} -gt ${portRangeEnd} ]; then
    echo "Given port to map out of range"
    exit 1
  fi
}

function store_mapping() {
  echo "store allocated port mapping: ${instance_id} -> $HOST_PORT"
  etcdctl set hivy/mapping/${instance_id} $HOST_PORT
}

# main() {
  echo
  echo "New member joined (event ${SERF_EVENT}). Parsing data..."
  while read line; do
      printf "Payload: ${line}\n"
      # name     ip     role
      # sample: xavier-local-machine-2  10.0.3.106      lab
      export instance_id=$(echo $line | awk '{print $1}')
      export VIRTUAL_IP=$(echo $line | awk '{print $2}')
      export ROLE=$(echo $line | awk '{print $3}')

      if [ "${ROLE}" != "lab" ]; then
        echo "${instance_id} is not a lab, ignoring member join."
        exit 0
      fi
  done

  export HOST_PORT=$(allocate_port)
  store_mapping

  echo "New lab ${instance_id} is ready, map ${HOST_PORT} port to ${VIRTUAL_IP}:22"
  map_virtual_port
#
