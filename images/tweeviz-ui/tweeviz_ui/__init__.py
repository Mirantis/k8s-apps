#!/usr/bin/env python

import os

import flask


header = os.getenv('TWEEVIZ_HEADER', 'Twitter stats')
tweeviz_api = os.environ["TWEEVIZ_API"]

app = flask.Flask(__name__)


@app.route('/')
@app.route('/index.html')
def index():
    return flask.render_template('index.html', tweeviz_api=tweeviz_api,
                                 header=header)


@app.route('/index2.html')
def index2():
    return flask.render_template('index2.html', tweeviz_api=tweeviz_api,
                                 header=header)


def main():
    app.run(host='0.0.0.0', port=8589)


if __name__ == "__main__":
    main()
