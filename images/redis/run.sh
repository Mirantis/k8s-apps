#!/bin/bash

set -e

function launchmaster {
  if [[ ! -e $MASTER_DATA ]]; then
    echo "Redis master data doesn't exist, data won't be persistent!"
    mkdir $MASTER_DATA
  fi
  master_conf=/opt/redis/conf/redis.conf
  sed -i "s|%data-dir%|$MASTER_DATA|" ${master_conf}
  redis-server ${master_conf} --protected-mode no
}

function launchsentinel {
  while true; do
    master=$(redis-cli -h ${!cluster_host} -p ${!cluster_port} --csv SENTINEL get-master-addr-by-name mymaster | tr ',' ' ' | cut -d' ' -f1)
    if [[ -n ${master} ]]; then
      master="${master//\"}"
    else
      master=$(hostname -i)
    fi
    if redis-cli -h ${master} INFO; then
      break
    fi
    echo "Connecting to master failed.  Waiting..."
    sleep 10
  done
  sentinel_conf=/opt/redis/conf/sentinel.conf
  sed -i "s/%master-ip%/${master}/" ${sentinel_conf}
  redis-sentinel ${sentinel_conf} --protected-mode no
}

function launchslave {
  while true; do
    master=$(redis-cli -h ${!cluster_host} -p ${!cluster_port} --csv SENTINEL get-master-addr-by-name mymaster | tr ',' ' ' | cut -d' ' -f1)
    master="${master//\"}"
    if redis-cli -h ${master} INFO; then
      break
    fi
    echo "Connecting to master failed.  Waiting..."
    sleep 10
  done
  data_dir="./"
  slave_conf=/opt/redis/conf/redis.conf
  sed -i "s|%data-dir%|${data_dir}|" ${slave_conf}
  sed -i "s/#//" ${slave_conf}
  sed -i "s/%master-ip%/${master}/" ${slave_conf}
  redis-server ${slave_conf} --protected-mode no
}

cluster_host=$(echo $CLUSTER_SERVICE_NAME | sed 's/-/_/g' | awk '{print toupper($0)}')_SERVICE_HOST
cluster_port=$(echo $CLUSTER_SERVICE_NAME | sed 's/-/_/g' | awk '{print toupper($0)}')_SERVICE_PORT
if [[ "${SENTINEL}" == "true" ]]; then
  launchsentinel
else
  mkdir -p /var/log/redis
  master=$(redis-cli -h ${!cluster_host} -p ${!cluster_port} --csv SENTINEL get-master-addr-by-name mymaster | tr ',' ' ' | cut -d' ' -f1)
  if [[ -n ${master} ]]; then
    launchslave
  else
    launchmaster
  fi
fi
