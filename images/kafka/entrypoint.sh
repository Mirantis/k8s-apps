#!/usr/bin/env bash

set -e

KAFKA_PORT=${KAFKA_PORT:-9092}
KAFKA_USER=${KAFKA_USER:-"kafka"}
KAFKA_DATA_DIR=${KAFKA_DATA_DIR:-"/var/lib/kafka/data"}

HOST=`hostname -s`
DOMAIN=`hostname -d`
MY_ID="${HOST##*-}"

function create_data_dirs {
    echo "Creating Kafka data directory and setting permissions"
    if [ ! -d $KAFKA_DATA_DIR  ]; then
        mkdir -p $KAFKA_DATA_DIR
        chown -R $KAFKA_USER:$KAFKA_USER $KAFKA_DATA_DIR
    fi
    echo "Created Kafka data directory and set permissions in $KAFKA_DATA_DIR"
}

create_data_dirs
/opt/kafka/bin/kafka-server-start.sh /opt/kafka/config/server.properties \
    --override broker.id=$MY_ID \
    --override zookeeper.connect=$ZK_CONNECT \
    --override advertised.listeners=PLAINTEXT://$HOST.$DOMAIN:$KAFKA_PORT \
    --override listeners=PLAINTEXT://:$KAFKA_PORT \
    --override log.dirs=$KAFKA_DATA_DIR

