FROM python:2.7-alpine

COPY . /tweeviz
RUN pip install /tweeviz

WORKDIR /tweeviz

ENTRYPOINT ["./tweeviz.py"]
