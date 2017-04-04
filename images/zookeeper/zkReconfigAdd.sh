#!/usr/bin/env bash

ZK_CLIENT_PORT=${ZK_CLIENT_PORT:-2181}
ZK_SERVER_PORT=${ZK_SERVER_PORT:-2888}
ZK_ELECTION_PORT=${ZK_ELECTION_PORT:-3888}

HOST=`hostname -s`
DOMAIN=`hostname -d`
MY_ID="${HOST##*-}"

echo "Adding myself to ensemble..."
while [[ $CHECK != 0 ]]; do
    /opt/zookeeper/bin/zkCli.sh reconfig -add "server.$MY_ID=$HOST.$DOMAIN:$ZK_SERVER_PORT:$ZK_ELECTION_PORT:participant;$ZK_CLIENT_PORT"
    CHECK=$?
    sleep 1
done
echo "Server was added to ensemble."
