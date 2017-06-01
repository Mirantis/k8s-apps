FROM openjdk:8 as builder

RUN git clone https://github.com/datastax/spark-cassandra-connector.git && cd spark-cassandra-connector && git checkout v2.0.2 && ./sbt/sbt -Dscala-2.11=true assembly


FROM openjdk:8 as compiler


ARG SBT_VERSION=0.13.13

WORKDIR /root/build
COPY --from=builder /spark-cassandra-connector/spark-cassandra-connector/target/full/scala-2.11/spark-cassandra-connector-assembly-2.0.2.jar lib/
COPY basic_test.scala build.sbt ./
RUN curl -L -o sbt-$SBT_VERSION.deb http://dl.bintray.com/sbt/debian/sbt-$SBT_VERSION.deb && \
    dpkg -i sbt-$SBT_VERSION.deb && \
    rm sbt-$SBT_VERSION.deb && \
    apt-get update && \
    apt-get install sbt && \
    sbt package


FROM openjdk:8

MAINTAINER Prashanth Babu <Prashanth.Babu@gmail.com>

# Scala related variables.
ARG SCALA_VERSION=2.11.8
ARG SCALA_BINARY_ARCHIVE_NAME=scala-${SCALA_VERSION}
ARG SCALA_BINARY_DOWNLOAD_URL=http://downloads.lightbend.com/scala/${SCALA_VERSION}/${SCALA_BINARY_ARCHIVE_NAME}.tgz

# Spark related variables.
ARG SPARK_VERSION=2.1.0
ARG SPARK_BINARY_ARCHIVE_NAME=spark-${SPARK_VERSION}-bin-hadoop2.7
ARG SPARK_BINARY_DOWNLOAD_URL=http://d3kbcqa49mib13.cloudfront.net/${SPARK_BINARY_ARCHIVE_NAME}.tgz

ARG JMX_VERSION=0.9
# Configure env variables for Scala, SBT and Spark.
# Also configure PATH env variable to include binary folders of Java, Scala, SBT and Spark.
ENV SCALA_HOME  /usr/local/scala
ENV SPARK_HOME  /usr/local/spark
ENV PATH        $JAVA_HOME/bin:$SCALA_HOME/bin:$SBT_HOME/bin:$SPARK_HOME/bin:$SPARK_HOME/sbin:$PATH

# Download, uncompress and move all the required packages and libraries to their corresponding directories in /usr/local/ folder.
RUN apt-get -yqq update && \
    apt-get install -yqq vim screen tmux && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /tmp/* && \
    wget -qO - ${SCALA_BINARY_DOWNLOAD_URL} | tar -xz -C /usr/local/ && \
    wget -qO - ${SPARK_BINARY_DOWNLOAD_URL} | tar -xz -C /usr/local/ && \
    cd /usr/local/ && \
    ln -s ${SCALA_BINARY_ARCHIVE_NAME} scala && \
    ln -s ${SPARK_BINARY_ARCHIVE_NAME} spark && \
    wget -P ${SPARK_HOME} https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/$JMX_VERSION/jmx_prometheus_javaagent-$JMX_VERSION.jar
COPY --from=builder /spark-cassandra-connector/spark-cassandra-connector/target/full/scala-2.11/spark-cassandra-connector-assembly-2.0.2.jar $SPARK_HOME/jars/
COPY --from=compiler /root/build/target/scala-2.11/simple-project_2.11-1.0.jar  $SPARK_HOME/jars/
