FROM openjdk:8

RUN curl -o /tmp/kibana.tgz https://artifacts.elastic.co/downloads/kibana/kibana-5.2.2-linux-x86_64.tar.gz \
    && tar -zxvf /tmp/kibana.tgz -C /usr/share/ \
    && mv /usr/share/kibana-5.2.2-linux-x86_64 /usr/share/kibana \
    && rm /tmp/kibana.tgz \
    && useradd --user-group --home-dir /usr/share/kibana kibana \
    && ln -s /usr/share/kibana/bin/kibana /usr/bin/kibana \
    && chown -R kibana:kibana /usr/share/kibana /usr/bin/kibana

USER kibana
