#!/usr/bin/env bash

set -e

HOST=`hostname -s`
MY_ID="${HOST##*-}"

/opt/zookeeper/bin/zkCli.sh reconfig -remove $MY_ID
