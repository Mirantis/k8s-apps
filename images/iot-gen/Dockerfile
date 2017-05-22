FROM golang:1.8.1

COPY src/main.go .
RUN go get -d -v "github.com/Masterminds/sprig" "github.com/Shopify/sarama" \
    "github.com/urfave/cli" "golang.org/x/time/rate"

RUN GOOS=linux go build -a -o iot-gen .

ENTRYPOINT ["./iot-gen"]

