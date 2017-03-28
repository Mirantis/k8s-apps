#!/usr/bin/env bash

set -e

ZK_CLIENT_PORT=${ZK_CLIENT_PORT:-2181}
ZK_SERVER_PORT=${ZK_SERVER_PORT:-2888}
ZK_ELECTION_PORT=${ZK_ELECTION_PORT:-3888}

HOST=`hostname -s`
DOMAIN=`hostname -d`
MY_ID="${HOST##*-}"

while [[ $CHECK != "imok" ]]; do
    CHECK=$(echo ruok | nc 127.0.0.1 $ZK_CLIENT_PORT)
    sleep 1
done

/opt/zookeeper/bin/zkCli.sh reconfig -add "server.$MY_ID=$HOST.$DOMAIN:$ZK_SERVER_PORT:$ZK_ELECTION_PORT:participant;$ZK_CLIENT_PORT"
