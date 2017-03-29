#!/usr/bin/env bash

set -e

ZK_REPLICAS=$1
ZK_SERVICE=$2
DOMAIN=`hostname -d`

while [[ $ZK_READY != $ZK_REPLICAS ]]; do
    ZK_READY=$(dig SRV ${DOMAIN/${DOMAIN%%.*}/$ZK_SERVICE} +short | wc -l)
    echo "Waiting for ZooKeeper start..."
    echo "Replicas ready: $ZK_READY/$ZK_REPLICAS"
    sleep 3
done
echo "All ZooKeeper replicas are ready."
