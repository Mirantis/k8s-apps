#!/usr/bin/env python3

import json
import os
import socket
import sys

from kazoo import client


def main():
    zk = client.KazooClient(os.environ["ZK_CONNECT"])
    zk.retry(zk.start)
    try:
        fqdn = socket.getfqdn()
        for broker_id in zk.retry(zk.get_children, "/brokers/ids"):
            data = zk.retry(zk.get, "/brokers/ids/" + broker_id)[0]
            if fqdn == json.loads(data.decode("UTF-8"))["host"]:
                sys.exit(0)
        sys.exit(1)
    finally:
        zk.retry(zk.stop)


if __name__ == "__main__":
    main()
