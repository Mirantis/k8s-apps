FROM openjdk:8

RUN set -x \
    && curl -o /tmp/logstash.tgz https://artifacts.elastic.co/downloads/logstash/logstash-5.2.2.tar.gz \
    && tar -zxvf /tmp/logstash.tgz -C /usr/share \
    && mv /usr/share/logstash-5.2.2 /usr/share/logstash \
    && rm /tmp/logstash.tgz \
    && useradd --user-group --home-dir /usr/share/logstash logstash \
    && ln -s /usr/share/logstash/bin/logstash /usr/bin/logstash \
    && chown -R logstash:logstash /usr/share/logstash /usr/bin/logstash

RUN /usr/share/logstash/bin/logstash-plugin update logstash-input-beats

USER logstash
