#!/bin/bash
set -e

echo "===Resolv dnsname==="
export DOCKER_VERNEMQ_DISCOVERY_NODE=`dig +short $DOCKER_VERNEMQ_DISCOVERY_NODE`
echo $DOCKER_VERNEMQ_DISCOVERY_NODE

start_vernemq
