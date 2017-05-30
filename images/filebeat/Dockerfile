FROM openjdk:8

RUN curl -o /tmp/filebeat.tgz https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-5.2.2-linux-x86_64.tar.gz \
    && tar -zxvf /tmp/filebeat.tgz -C /usr/share \
    && mv /usr/share/filebeat-5.2.2-linux-x86_64 /usr/share/filebeat \
    && rm /tmp/filebeat.tgz \
    && useradd --user-group --home-dir /usr/share/filebeat filebeat \
    && ln -s /usr/share/filebeat/filebeat /usr/bin/filebeat \
    && mkdir -p /etc/filebeat \
    && chown -R filebeat:filebeat /usr/share/filebeat /usr/bin/filebeat /etc/filebeat
