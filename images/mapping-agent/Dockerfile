FROM quay.io/deis/base:v0.3.7

RUN curl -s -L https://github.com/dcos/demos/releases/download/0.2/mapping-agent -o /mapping-agent && \
    chmod +x /mapping-agent && \
    mkdir /content && \
    curl -s -L https://raw.githubusercontent.com/dcos/demos/master/sensoranalytics/1.9/mapping-agent/content/index.html -o /content/index.html && \
    curl -s -L https://raw.githubusercontent.com/dcos/demos/master/sensoranalytics/1.9/mapping-agent/content/marker.png -o /content/marker.png

ENTRYPOINT ["/mapping-agent"]