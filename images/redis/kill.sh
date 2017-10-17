#!/bin/bash

PET_ORDINAL="${HOSTNAME##*-}"
MASTER_ADDRESS="${HOSTNAME%-*}-0.${HOSTNAME%-*}"
PEER_IP=$(perl -MSocket -e 'print inet_ntoa(scalar(gethostbyname($ARGV[0])))' "${MASTER_ADDRESS}")

NODE_ID="$(redis-cli cluster nodes | grep myself | cut -d" " -f1)"
ruby /usr/local/bin/redis-trib.rb call $PEER_IP:6379 CLUSTER FORGET $NODE_ID
yes yes | ruby /usr/local/bin/redis-trib.rb fix $PEER_IP:6379
