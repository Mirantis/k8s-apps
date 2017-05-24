#!/usr/bin/env python

import json
import sys

import click
import kafka as kafka_client
import tweepy


def process_comma_separated_option(option):
    return option.split(',') if option else None


@click.command()
@click.option('--app_key', help="Twitter API Application Key", required=True)
@click.option('--app_secret', help="Twitter API Application Secret",
              required=True)
@click.option('--token_key', help="Twitter API Token Key", required=True)
@click.option('--token_secret', help="Twitter API Token Secret", required=True)
@click.option('--follow', '-f', default="",
              help="A comma-separated list of user IDs to follow")
@click.option('--track', '-t', default="",
              help="A comma-separated list of phrases to follow")
@click.option('--locations', '-l', default="",
              help="A comma-separated list of longitude,latitude pairs"
                   " specifying a set of bounding boxes for locations to"
                   " follow Tweets from")
@click.option('--languages', default='en',
              help="Limit filter to the comma-separated list of BCP 47"
                   " language identifiers")
@click.option('--kafka', default="localhost:9092",
              help="A comma-separated list of Kafka bootstap servers")
@click.option('--topic', default="twitter-stream",
              help="Kafka topic where the tweets will be published")
def _main(app_key, app_secret, token_key, token_secret, follow, track,
         locations, languages, kafka, topic):
    """TweePub reads tweets from Twitter Streaming API with provided
       characteristics and pushes them to specified Apache Kafka instance.
    """
    follow = follow.split(',') if follow else None
    track = track.split(',') if track else None
    locations = [float(l) for l in locations.split(',')] if locations else None
    languages = languages.split(',') if languages else None

    if not any([follow, track, locations]):
        click.echo("Error: At least one of follow, track or locations should"
                   " be specified")
        sys.exit(-1)

    if follow: click.echo("Follow: %s" % follow)
    if track: click.echo("Track: %s" % track)
    if locations: click.echo("Locations: %s" % locations)
    if languages: click.echo("Languages: %s" % languages)

    # Auth in Twitter and check if it's succeded
    auth = tweepy.OAuthHandler(app_key, app_secret)
    auth.set_access_token(token_key, token_secret)
    api = tweepy.API(auth)
    click.echo("Twitter authenticated user: %s" % api.me().screen_name)

    click.echo("Kafka bootstrap servers: %s" % kafka)
    producer = kafka_client.KafkaProducer(bootstrap_servers=kafka,
                                          value_serializer=str.encode)

    stream_listener = TwitterStreamListener(producer, topic)
    stream = tweepy.Stream(auth=api.auth, listener=stream_listener)
    stream.filter(follow=follow, track=track, locations=locations, languages=languages)
    producer.flush()


class TwitterStreamListener(tweepy.StreamListener):
    def __init__(self, producer, topic):
        self.producer = producer
        self.topic = topic
        super(TwitterStreamListener, self).__init__()

    def on_status(self, status):
        click.echo(">>> %s ::: %s" % (status.author.screen_name, status.text))
        self.producer.send(self.topic, json.dumps(status._json))


    def on_error(self, status_code):
        click.echo("Error happened in Twitter Stream listener: %s" % status_code)
        # Disconnect only if rate limited and keep going otherwise
        if status_code == 420:
            return False


def main():
    _main(auto_envvar_prefix='TWEEPUB')

if __name__ == '__main__':
    main()
