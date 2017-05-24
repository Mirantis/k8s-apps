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
stats = {'popularity': [], 'top': []}
processed_results = set()


def update_stats():
    hdfs = client.Client(hdfs_address, hdfs_port, use_trash=False)

    parts = []
    for result in sorted([r['path'] for r in hdfs.ls([results_dir])]):
        if not hdfs.test(result + "/_SUCCESS", exists=True):
            print "No _SUCCESS in %s" % result
        if result in processed_results:
            continue
        processed_results.add(result)
        parts += sorted([r['path'] for r in hdfs.ls([result + "/part*"])])

    if not parts:
        return

    print "New data in: %s" % parts

    for part in hdfs.text(parts):
        part_stats = eval("[" + ",".join(part.split('\n')) + "]")
        for stat in part_stats:
            hashtag = stat[0].lower()
            hashtags[hashtag] = hashtags.get(hashtag, 0) + stat[1][0]

    stats['popularity'] = to_jqcloud_format(filter(lambda x: x[1] >= min_popularity, hashtags.items()))
    max_top = min(len(hashtags), top_list_len)
    if not max_top:
        return
    stats['top'] = to_jqcloud_format(sorted(hashtags.items(), key=lambda x: x[1], reverse=True)[:max_top])


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
