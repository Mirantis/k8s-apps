#!/bin/bash

set -ex

function check_status {
    if ruby /usr/local/bin/redis-trib.rb check $PEER_IP:6379 | grep '^\[ERR\]'; then
        return 1
    else
        return 0
    fi
}

function join_cluster {
    PET_ORDINAL="${HOSTNAME##*-}"
    MASTER_ADDRESS="${HOSTNAME%-*}-0.${HOSTNAME%-*}"

    if [ $PET_ORDINAL = "0" ]; then
        # The first member of the cluster should control all slots initially
        echo "Bootstrapping this cluster node with all cluster slots..."
        redis-cli cluster addslots $(seq 0 16383)
    else
        # TODO: Get list of peers using the peer finder using an init container
        PEER_IP=$(perl -MSocket -e 'print inet_ntoa(scalar(gethostbyname($ARGV[0])))' "${MASTER_ADDRESS}")
        # TODO: Make sure the node we're initializing is not already a master (it may be a recovering node)
        redis-cli cluster meet $PEER_IP 6379
        # check status
        until check_status; do
            sleep 5
        done
        # rebalance slots
        ruby /usr/local/bin/redis-trib.rb rebalance --auto-weights --use-empty-masters $PEER_IP:6379
    fi
}

mkdir -p /var/log/redis

redis-server /opt/redis/conf/redis.conf &

# TODO: Wait until redis-server process is ready
sleep 1

if [[ "${1}" == "cluster" ]]; then
    join_cluster
fi

wait
