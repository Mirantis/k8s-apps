#!/usr/bin/env python

import json
import os
import threading
import time

import flask
from snakebite import client


hdfs_address = os.getenv('TWEEVIZ_HDFS_ADDRESS', 'hdfs-namenode')
hdfs_port = int(os.getenv('TWEEVIZ_HDFS_PORT', 8020))
results_dir = os.getenv('TWEEVIZ_HDFS_PATH', '/')
min_popularity = int(os.getenv('TWEEVIZ_MIN_POPULARITY', 2))
top_list_len = int(os.getenv('TWEEVIZ_TOP_LIST_SIZE', 0))


hashtags = {}
stats = {'popularity': []}
processed_results = set()


def update_stats():
    hdfs = client.Client(hdfs_address, hdfs_port, use_trash=False)

    parts = []
    for result in sorted([r['path'] for r in hdfs.ls([results_dir])]):
        if not hdfs.test(result + "/_SUCCESS", exists=True):
            continue
        if result in processed_results:
            continue
        processed_results.add(result)
        parts += sorted([r['path'] for r in hdfs.ls([result + "/part*"])])

    if not parts:
        return

    for part in hdfs.text(parts):
        part_stats = part.split('\n')
        for stat in part_stats:
            if not stat:
                continue
            hashtag, count = stat.split(' ')
            hashtags[hashtag] = hashtags.get(hashtag, 0) + int(count)

    sorted_stats = to_jqcloud_format(sorted(hashtags.items(), key=lambda x: x[1], reverse=True))
    max_top = min(len(hashtags), top_list_len)
    stats['popularity'] = sorted_stats if top_list_len == 0 else sorted_stats[:max_top]

    print "Processed data in: %s" % parts


def to_jqcloud_format(keypairs):
    return [{
        'text': kp[0],
        'weight': kp[1],
    } for kp in keypairs]


app = flask.Flask(__name__)


@app.route('/')
@app.route('/index.html')
def index():
    return flask.render_template('index.html')


@app.route('/index2.html')
def index2():
    return flask.render_template('index2.html')


@app.route('/stats')
def get_stats():
    response = flask.jsonify(stats)
    response.headers.add('Access-Control-Allow-Origin', '*')
    return response


def stats_updater():
    while(True):
        try:
            update_stats()
        except Exception:
            pass
        time.sleep(1)


def main():
    thread = threading.Thread(target=stats_updater)
    thread.daemon = True
    thread.start()

    app.run(host='0.0.0.0', port=8589)


if __name__ == "__main__":
    main()
