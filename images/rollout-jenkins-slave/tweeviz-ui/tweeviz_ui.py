#!/usr/bin/env python

from multiprocessing import Value
import os

import flask


header = os.getenv('TWEEVIZ_HEADER', 'Twitter stats')
tweeviz_api = os.environ["TWEEVIZ_API"]

requests_counter = Value('i', 0)

app = flask.Flask(__name__)


@app.route('/')
@app.route('/index.html')
def index():
    with requests_counter.get_lock():
        requests_counter.value += 1

#    if requests_counter.value > 50:
#        flask.abort(500)

    return flask.render_template('index.html', tweeviz_api=tweeviz_api,
                                 header=header)

def main():
    app.run(host='0.0.0.0', port=8589)


if __name__ == "__main__":
    main()
