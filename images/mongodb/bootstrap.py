import json
import logging
import os
import socket
import subprocess
import sys
import time

import bson
import bson.objectid
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


def _start_mongo(command):
    params = {
        "shell": True,
        "stdin": sys.stdin,
        "stdout": sys.stdout,
        "stderr": sys.stderr
    }
    return subprocess.Popen(command, **params)


def start_shard(port):
    command = ("mongod --logpath /var/log/mongodb/mongo.log --shardsvr"
               " --replSet shard --port %d") % port
    return _start_mongo(command)


def start_configsvr(port):
    command = ("mongod --logpath /var/log/mongodb/mongo.log --configsvr"
               " --replSet config --port %d") % port
    return _start_mongo(command)


def start_router(port, configdb_address):
    command = ("mongos --logpath /var/log/mongodb/mongo.log"
               " --configdb config/%s --port %d") % (configdb_address, port)
    return _start_mongo(command)


def init_replicaset(master_address, replicaset):
    LOG.info("Init replicaset")
    client = pymongo.MongoClient(host=[master_address])
    init_config = {
        "_id": replicaset,
        "members": [
            {"_id": 0, "host": master_address}
        ]
    }
    res = client.admin.command("replSetInitiate", init_config)
    LOG.debug("Replicaset was initialized, status:\n%s" % prettyjson(res))


def join_to_replicaset(master_address, host_id, port, replicaset,
                       is_configsvr=False):
    client = pymongo.MongoClient(host=[master_address],
                                 replicaset=replicaset)
    current_config = client.admin.command("replSetGetConfig")
    new_config = {
        "_id": replicaset,
        "configsvr": is_configsvr,
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
            break
        except Exception:
            LOG.exception("Reconfig failed, retrying...")
            time.sleep(0.2)


def add_shard(router_address, master_address):
    while True:
        try:
            client = pymongo.MongoClient(host=[router_address])
            res = client.admin.command("addShard", "shard/%s" % master_address)
            LOG.debug("Shard status:\n%s" % prettyjson(res))
            break
        except Exception:
            LOG.exception("Adding shard failed, retrying...")
            time.sleep(0.2)


def main():
    port = int(sys.argv[1])
    replicaset = sys.argv[2]
    if replicaset == "config":
        daemon = start_configsvr(port)
        is_configsvr = True
    elif replicaset == "shard":
        daemon = start_shard(port)
        is_configsvr = False
    elif replicaset == "router":
        config_address = os.environ.get("MONGO_CONFIGDB_ADDRESS")
        if config_address is None:
            LOG.error("Config DB address is not specified"
                      "in MONGO_CONFIGDB_ADDRESS environment variable")
            sys.exit(1)
        daemon = start_router(port, config_address)
    if replicaset != "router":
        wait_local_mongo(port)
        master_address = get_master_address(port)
        host_id = get_host_id()
        if host_id == 0:
            init_replicaset(master_address, replicaset)
            if replicaset == "shard":
                router_address = os.environ.get("MONGO_ROUTER_ADDRESS")
                if router_address is None:
                    LOG.error("Router address is not specified in "
                              "MONGO_ROUTER_ADDRESS environment variable")
                    sys.exit(1)
                add_shard(router_address, master_address)
        else:
            join_to_replicaset(master_address, host_id, port, replicaset,
                            is_configsvr)
    daemon.wait()
    daemon.communicate()


if __name__ == "__main__":
    main()
