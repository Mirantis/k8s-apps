#!/usr/bin/env python

import json
import sys

import pyspark
from pyspark import streaming
from pyspark.streaming import kafka
from pyspark.sql import SQLContext, Row
from pyspark.sql.types import *

def get_hashtags(tweet):
    data = json.loads(tweet)
    data_hashtags = data.get("entities", {}).get("hashtags", [])
    hashtags = []

    for hashtag in data_hashtags:
        # filter out unicode
        try:
            hashtags.append(str("#" + hashtag["text"]).lower())
        except (UnicodeEncodeError, KeyError):
            pass

    return hashtags

def create_writer(keyspace, mode="append"):
    def writer(df, table):
        df.write.format("org.apache.spark.sql.cassandra").\
                 options(table=table, keyspace=keyspace).save(mode="append")
    return writer

def getSqlContextInstance(sparkContext):
    if ('sqlContextSingletonInstance' not in globals()):
        globals()['sqlContextSingletonInstance'] = SQLContext(sparkContext)
    return globals()['sqlContextSingletonInstance']

def getSqlWriter(keyspace):
    if ('sqlWriter' not in globals()):
      globals()['sqlWriter'] = create_writer(keyspace)
    return globals()['sqlWriter']

def write_to_cas(time, rdd):
    print "============== %s ============" % str(time)
    _, keyspace, table = storage.split(":")
    local_sql = getSqlContextInstance(rdd.context)
    data = local_sql.read.format("org.apache.spark.sql.cassandra").load(keyspace=keyspace, table=table)
    writer = getSqlWriter(keyspace)
    rdd = rdd.map(lambda line: line.split())
    rdd = rdd.map(lambda(_, count): (_, int(count)))
    schema = StructType([
      StructField("hashtag", StringType(), False),
      StructField("count", IntegerType(), False)])
    stream = rdd.toDF(schema)
    res = data.union(stream)
    from operator import add
    result = local_sql.createDataFrame(res.rdd.foldByKey(0, add), schema)
    writer(result, table)

    print "========== DONE WRITING ============== "

def create_keyspace(storage):
    from cassandra.cluster import Cluster

    cluster = Cluster([storage.split(":")[0]])
    session = cluster.connect()

    session.execute("create KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3}" % (storage.split(":")[1]))
    session.execute("create table IF NOT EXISTS {}.{} (hashtag text primary key, count int)".format(storage.split(":")[1],storage.split(":")[2]))


def save(stream, save_to, storage):
    if 'cassandra' in save_to:
      create_keyspace(storage)
      globals()['storage'] = storage
      stream.foreachRDD(write_to_cas)
    else:
      stream.saveAsTextFiles(storage)

def main():
    if len(sys.argv) == 2 and sys.argv[1] == "noop":
        return
    if len(sys.argv) != 8:
        print "Usage: spark_hashtags_count.py <spark_master> <zk_quorum> <topic_name> <min_hashtag_counts> <batch_duration> <save_to>"
        print "Example: spark_hashtags_count.py local[4] zk-kafka-1-0.zk-kafka-1:2181,zk-kafka-1-1.zk-kafka-1:2181,zk-kafka-1-2.zk-kafka-1:2181 twitter-stream 0 5 hdfs://hdfs-namenode:8020/demo"
        print "<spark_master> - spark master to use: local[4] or spark://HOST:PORT"
        print "<zk_quorum> - zk quorum to connect: zk-kafka-1-0.zk-kafka-1:2181,zk-kafka-1-1.zk-kafka-1:2181,zk-kafka-1-2.zk-kafka-1:2181"
        print "<topic_name> - kafka topic name: twitter-stream"
        print "<min_hashtag_counts> - filter out hashtags with less then specified count"
        print "<batch_duration> - spark streaming batch duration ~ how often data will be written"
        print "<save_to> - hdfs or cassandra"
        print "<storage> - save as text files to: hdfs://hdfs-namenode:8020/demo or to database: <host>:<keyspace>:<table>"
        exit(-1)

    spark_master = sys.argv[1]
    zk_quorum = sys.argv[2]
    topic_name = sys.argv[3]
    min_hashtag_counts = int(sys.argv[4])
    batch_duration = int(sys.argv[5])
    save_to = sys.argv[6]
    storage = sys.argv[7]
    sc = pyspark.SparkContext(spark_master, appName="TweeTics")
    ssc = streaming.StreamingContext(sc, batch_duration)
    sql = SQLContext(sc)

    tweets = kafka.KafkaUtils.createStream(ssc, zk_quorum, "tweetics-consumer", {topic_name: 1}).map(lambda x: x[1])
    counts = tweets.flatMap(get_hashtags).map(lambda hashtag: (hashtag, 1)).reduceByKey(lambda a, b: a + b)
    sorted_counts = counts.transform(lambda rdd: rdd.sortByKey(ascending=False, keyfunc=lambda x: x[1]))
    output = sorted_counts.map(lambda x: "%s %s" % (x[0], x[1]))

    output.pprint()
    save(output, save_to, storage)

    ssc.start()
    ssc.awaitTermination()


if __name__ == "__main__":
    main()
