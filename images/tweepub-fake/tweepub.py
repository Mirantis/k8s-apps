#!/usr/bin/env python

import json
import random
import time

import click
import kafka as kafka_client


@click.command()
@click.option('--base-word', '-w', default="hashtag",
              help="Base word to be used in hashtags of fake tweets")
@click.option('--kafka', default=None,
              help="A comma-separated list of Kafka bootstrap servers")
@click.option('--topic', default="twitter-stream",
              help="Kafka topic where the fake tweets will be published")
def _main(base_word, kafka, topic):
    """TweePub reads tweets from Twitter Streaming API with provided
       characteristics and pushes them to specified Apache Kafka instance.
    """
    click.echo("Using base word: %s" % base_word)
    click.echo("Kafka bootstrap servers: %s" % kafka)
    producer = kafka_client.KafkaProducer(bootstrap_servers=kafka,
                                          value_serializer=str.encode)

    for _ in range(0, 100):
        hashtags = []

        for i in range(1, 16):
            if random.randint(1, i+1) == i:
                hashtags.append(base_word + "-" + str(i))

        click.echo("Pushing hashtags: %s" % hashtags)

        producer.send(topic, json.dumps(dict(
            entities=dict(
                hashtags=[
                    dict(text=ht) for ht in hashtags
                ]
            )
        )))
        time.sleep(5)

    producer.flush()


def main():
    _main(auto_envvar_prefix='TWEEPUB')


if __name__ == '__main__':
    main()
