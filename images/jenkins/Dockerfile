FROM jenkins/jenkins:2.104
RUN \
  install-plugins.sh kubernetes:1.1.4 workflow-aggregator:2.5 workflow-job:2.17 credentials-binding:1.14 git:3.7.0 && \
  mv /usr/share/jenkins/ref/plugins /usr/share/jenkins/ref/bundled_plugins
COPY install-plugins.sh /usr/local/bin/install-plugins.sh
