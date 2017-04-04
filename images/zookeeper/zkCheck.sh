#!/usr/bin/env bash

set -e
while [[ $CHECK != "imok" ]]; do
    CHECK=$(echo ruok | nc 127.0.0.1 $ZK_CLIENT_PORT)
    sleep 1
done
