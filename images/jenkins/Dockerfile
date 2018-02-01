FROM jenkins/jenkins:2.88
RUN \
  install-plugins.sh kubernetes:0.11 workflow-aggregator:2.5 workflow-job:2.13 credentials-binding:1.12 git:3.4.0 && \
  mv /usr/share/jenkins/ref/plugins /usr/share/jenkins/ref/bundled_plugins
COPY install-plugins.sh /usr/local/bin/install-plugins.sh
