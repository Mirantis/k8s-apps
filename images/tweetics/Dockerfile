FROM mirantisworkloads/spark:2.1.1

ADD spark_hashtags_count.py /
RUN git clone https://github.com/datastax/python-driver.git && cd python-driver && python setup.py install && cd /
RUN spark-submit --packages org.apache.spark:spark-streaming-kafka-0-8_2.11:2.1.0 spark_hashtags_count.py noop
