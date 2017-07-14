from __future__ import print_function

import os
import sys

import pymongo


def msg(service_name):
    def msg_decorator(fct):
        def fct_wrapper(*args, **kwargs):
            print(">>> Check %s service..." % service_name)
            res = fct(*args, **kwargs)
            print(">>> Done %s" % service_name)
            return res
        return fct_wrapper
    return msg_decorator


def _check_replicaset(address, replicas, replicaset):
    success = True
    client = pymongo.MongoClient(host=address,
                                 replicaset=replicaset)
    config = client.admin.command("replSetGetStatus")
    # check status
    print("Replicaset status: %s" % config["ok"])
    success = success and (config["ok"] == 1)
    # check replicas
    print("Replicas count: %d" % len(config["members"]))
    success = success and (len(config["members"]) == replicas)
    # check nodes
    for member in config["members"]:
        print("Status '%s': %s" % (member["name"], member["health"]))
        success = success and (member["health"] == 1)
    client.close()
    return success


@msg("Mongo Config")
def check_cfg(address, replicas):
    return _check_replicaset(address, replicas, "config")


@msg("Mongo Shard")
def check_shard(address, replicas):
    return _check_replicaset(address, replicas, "shard")


@msg("Mongo router")
def check_router(address, replicas, shard_replicas):
    success = True
    client = pymongo.MongoClient(host=address)
    status = client.admin.command("serverStatus")
    # check cluster status
    print("Cluster status: %s" % status["ok"])
    success = success and (status["ok"] == 1)
    return success


def main():
    cfg_address = os.environ["MONGO_CONFIGDB_ADDRESS"].split(",")
    shard_address = os.environ["MONGO_SHARD_ADDRESS"].split(",")
    router_address = [os.environ["MONGO_ROUTER_ADDRESS"]]

    cfg_replicas = int(os.environ["MONGO_CONFIGDB_REPLICAS"])
    shard_replicas = int(os.environ["MONGO_SHARD_REPLICAS"])
    router_replicas = int(os.environ["MONGO_ROUTER_REPLICAS"])

    status = all((
        check_cfg(cfg_address, cfg_replicas),
        check_shard(shard_address, shard_replicas),
        check_router(router_address, router_replicas, shard_replicas)))
    if not status:
        sys.exit(1)


if __name__ == "__main__":
    main()
