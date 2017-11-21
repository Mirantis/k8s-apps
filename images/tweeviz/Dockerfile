FROM python:2.7-jessie

COPY . /tweeviz
RUN pip install /tweeviz

WORKDIR /tweeviz

ENTRYPOINT ["./tweeviz.py"]
