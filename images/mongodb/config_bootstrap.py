import json
import logging
import socket
import subprocess
import sys
import time

import bson
import pymongo


def get_logger():
    logger = logging.getLogger("mongo_bootstrap")
    handler = logging.FileHandler("/mongo_bootstrap.log")
    formatter = logging.Formatter("%(asctime)s %(levelname)s: %(message)s")
    handler.setFormatter(formatter)
    logger.addHandler(handler)
    logger.setLevel(logging.DEBUG)
    return logger


LOG = get_logger()


class BSONEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, (bson.ObjectId, bson.Timestamp)):
            return str(obj)
        return json.JSONEncoder.default(self, obj)


def prettyjson(obj):
    return json.dumps(obj, sort_keys=True, indent=4, cls=BSONEncoder)


def get_master_address(port):
    hostname = socket.gethostname()
    service_name, _, host_id = hostname.rpartition("-")
    if host_id == "0":
        return get_local_address(port)
    fqdn = socket.getfqdn()
    domain = fqdn[len(hostname) + 1:]
    return "%(service_name)s-0.%(domain)s:%(port)s" % {
        "service_name": service_name,
        "domain": domain,
        "port": port}


def get_local_address(port):
    return "%s:%s" % (socket.getfqdn(), port)


def get_host_id():
    _, _, host_id = socket.gethostname().rpartition("-")
    return int(host_id)


def wait_local_mongo(port):
    client = pymongo.MongoClient(host=[get_local_address(port)])
    while True:
        try:
            info = client.server_info()
            LOG.debug("Local mongo status:\n%s" % prettyjson(info))
            return
        except Exception:
            LOG.exception("Local mongo is not available")
            time.sleep(0.2)


def start_mongo(port):
    params = {
        "shell": True,
        "stdin": sys.stdin,
        "stdout": sys.stdout,
        "stderr": sys.stderr
    }
    command = "mongod --configsvr --replSet config --port %d" % port
    return subprocess.Popen(command, **params)


def init_replicaset(master_address):
    LOG.info("Init replicaset")
    client = pymongo.MongoClient(host=[master_address])
    init_config = {
        "_id": "config",
        "members": [
            {"_id": 0, "host": master_address}
        ]
    }
    res = client.admin.command("replSetInitiate", init_config)
    LOG.debug("Replicaset was initialized, status:\n%s" % prettyjson(res))


def join_to_replicaset(master_address, host_id, port):
    client = pymongo.MongoClient(host=[master_address],
                                    replicaset="config")
    current_config = client.admin.command("replSetGetConfig")
    new_config = {
        "_id": "config",
        "configsvr": True,
        "version": current_config["config"]["version"] + 1,
        "members": []
    }
    # keep old members
    for member in current_config["config"]["members"]:
        new_config["members"].append({
            "_id": member["_id"],
            "host": member["host"]
        })
    # add new one
    new_config["members"].append({
        "_id": host_id,
        "host": get_local_address(port)
    })
    # reconfig
    while True:
        try:
            res = client.admin.command("replSetReconfig", new_config)
            LOG.debug("Replicaset status:\n%s" % prettyjson(res))
        except Exception:
            LOG.exception("Reconfig failed, retrying...")
            time.sleep(0.2)


def main():
    port = int(sys.argv[1])
    daemon = start_mongo(port)
    wait_local_mongo(port)
    master_address = get_master_address(port)
    host_id = get_host_id()
    if host_id == 0:
        init_replicaset(master_address)
    else:
        join_to_replicaset(master_address, host_id, port)
    daemon.wait()
    daemon.communicate()


if __name__ == "__main__":
    main()
