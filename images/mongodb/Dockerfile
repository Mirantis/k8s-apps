FROM mongo:3.4

COPY config_bootstrap.py /config_bootstrap.py

RUN apt update && \
    apt install -y python python-pip build-essential python-dev && \
    pip install pymongo && \
    apt purge -y build-essential python-dev
