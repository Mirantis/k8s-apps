#!/usr/bin/env bash

set -e

ZK_USER=${ZK_USER:-"zookeeper"}
ZK_DATA_DIR=${ZK_DATA_DIR:-"/var/lib/zookeeper/data"}
ZK_DATA_LOG_DIR=${ZK_DATA_LOG_DIR:-"/var/lib/zookeeper/log"}
ZK_LOG_DIR=${ZK_LOG_DIR:-"var/log/zookeeper"}
ZK_CONF_DIR=${ZK_CONF_DIR:-"/opt/zookeeper/conf"}
ZK_CONFIG_FILE="$ZK_CONF_DIR/zoo.cfg"
ZK_CONFIG_DYNAMIC="$ZK_CONF_DIR/zoo.cfg.dynamic"
ID_FILE="$ZK_DATA_DIR/myid"

ZK_CLIENT_PORT=${ZK_CLIENT_PORT:-2181}
ZK_SERVER_PORT=${ZK_SERVER_PORT:-2888}
ZK_ELECTION_PORT=${ZK_ELECTION_PORT:-3888}

HOST=`hostname -s`
DOMAIN=`hostname -d`
MY_ID="${HOST##*-}"

function create_data_dirs {
    echo "Creating ZooKeeper data directories and setting permissions"
    if [ ! -d $ZK_DATA_DIR  ]; then
        mkdir -p $ZK_DATA_DIR
        chown -R $ZK_USER:$ZK_USER $ZK_DATA_DIR
    fi

    if [ ! -d $ZK_DATA_LOG_DIR  ]; then
        mkdir -p $ZK_DATA_LOG_DIR
        chown -R $ZK_USER:$ZK_USER $ZK_DATA_LOG_DIR
    fi

    if [ ! -d $ZK_LOG_DIR  ]; then
        mkdir -p $ZK_LOG_DIR
        chown -R $ZK_USER:$ZK_USER $ZK_LOG_DIR
    fi

    if [ ! -f $ID_FILE ]; then
        echo $MY_ID > $ID_FILE
    fi
    echo "Created ZooKeeper data directories and set permissions in $ZK_DATA_DIR"
}

function create_config {
    echo "Creating Zookeeper static configuration"
    cp $ZK_CONF_DIR/static/zoo.cfg $ZK_CONF_DIR
    echo "Wrote ZooKeeper static configuration to $ZK_CONFIG_FILE"
    echo "Creating ZooKeeper dynamic configuration"
    echo "server.$MY_ID=$HOST.$DOMAIN:$ZK_SERVER_PORT:$ZK_ELECTION_PORT:participant;$ZK_CLIENT_PORT" >> $ZK_CONFIG_DYNAMIC
    dig SRV $DOMAIN +short | while read line; do
        fqdn="${line##* }"
        name="${fqdn%%.*}"
        id="${name##*-}"
        echo "server.$id=${fqdn%.}:$ZK_SERVER_PORT:$ZK_ELECTION_PORT:participant;$ZK_CLIENT_PORT" >> $ZK_CONFIG_DYNAMIC
    done
    echo "Wrote ZooKeeper dynamic configuration to $ZK_CONFIG_DYNAMIC"
}

create_config
create_data_dirs
/opt/zookeeper/bin/zkServer.sh start-foreground
