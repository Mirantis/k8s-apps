#!/usr/bin/env python

import json
import sys

import pyspark
from pyspark import streaming
from pyspark.streaming import kafka


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


if __name__ == "__main__":
    if len(sys.argv) != 7:
        print "Usage: spark_hashtags_count.py <spark_master> <zk_quorum> <topic_name> <min_hashtag_counts> <batch_duration> <save_to>"
        print "Example: spark_hashtags_count.py local[4] zk-kafka-1-0.zk-kafka-1:2181,zk-kafka-1-1.zk-kafka-1:2181,zk-kafka-1-2.zk-kafka-1:2181 twitter-stream 0 5 hdfs://hdfs-namenode:8020/demo"
        print "<spark_master> - spark master to use: local[4] or spark://HOST:PORT"
        print "<zk_quorum> - zk quorum to connect: zk-kafka-1-0.zk-kafka-1:2181,zk-kafka-1-1.zk-kafka-1:2181,zk-kafka-1-2.zk-kafka-1:2181"
        print "<topic_name> - kafka topic name: twitter-stream"
        print "<min_hashtag_counts> - filter out hashtags with less then specified count"
        print "<batch_duration> - spark streaming batch duration ~ how often data will be written"
        print "<save_to> - save as text files to: hdfs://hdfs-namenode:8020/demo"
        exit(-1)

    spark_master = sys.argv[1]
    zk_quorum = sys.argv[2]
    topic_name = sys.argv[3]
    min_hashtag_counts = int(sys.argv[4])
    batch_duration = int(sys.argv[5])
    save_to = sys.argv[6]

    sc = pyspark.SparkContext("local[2]", appName="TweeTics")
    ssc = streaming.StreamingContext(sc, batch_duration)

    tweets = kafka.KafkaUtils.createStream(ssc, zk_quorum, "tweetics-consumer", {topic_name: 1}).map(lambda x: x[1])
    counts = tweets.flatMap(get_hashtags).map(lambda hashtag: (hashtag, 1)).reduceByKey(lambda a, b: a + b)
    sorted_counts = counts.transform(lambda rdd: rdd.sortByKey(ascending=False, keyfunc=lambda x: x[1]))
    output = sorted_counts.map(lambda x: "%s %s" % (x[0], x[1]))

    output.pprint()
    output.saveAsTextFiles(save_to)

    ssc.start()
    ssc.awaitTermination()
