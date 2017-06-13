FROM openjdk:8

ARG ZEPPELIN_VERSION=0.7.2

# install zeppelin
RUN curl http://www-eu.apache.org/dist/zeppelin/zeppelin-${ZEPPELIN_VERSION}/zeppelin-${ZEPPELIN_VERSION}-bin-all.tgz -o zeppelin.tgz \
    && tar xvf zeppelin.tgz \
    && ln -s zeppelin-${ZEPPELIN_VERSION}-bin-all zeppelin\
    && rm zeppelin.tgz

# install matplotlib integration
RUN apt update \
    && apt install -y python-matplotlib

# install R integration
# FIXME: install r-cran-knitr from deb once it is fixed in jessie-backports
RUN apt install -y -t jessie-backports r-base && \
    R -e 'install.packages("knitr",repos = "http://cran.r-project.org")'

ENV PATH /zeppelin/bin:$PATH
