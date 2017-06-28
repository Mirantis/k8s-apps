FROM golang:1.8.1

COPY src/main.go .
RUN set -x && \
    git clone https://github.com/kubernetes/api $GOPATH/src/k8s.io/api && \
    go get -d -v "github.com/kelseyhightower/envconfig" "github.com/prometheus/client_golang/api" \
    "github.com/prometheus/common/model" "golang.org/x/time/rate" \
    "k8s.io/apimachinery/pkg/apis/meta/v1" "k8s.io/client-go/kubernetes" \
    "k8s.io/client-go/rest"


RUN GOOS=linux go build -a -o autoscale .

ENTRYPOINT ["./autoscale"]

