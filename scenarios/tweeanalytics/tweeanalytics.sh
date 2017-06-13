#!/bin/bash

set -ex

HELM_CMD=${HELM_CMD:-helm}
KUBECTL_CMD=${KUBECTL_CMD:-kubectl}
MODE="multi"
NAMESPACE=$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1)
TWITTER_SOCKET="tweeviz-tweeviz-$NAMESPACE.$NAMESPACE.svc.cluster.local:8589"
STORAGE="hdfs"

EXIT_CODE=0

for ((i=1;i<=$#;i++));
do
    case ${!i} in
        "--host")
        ((i++))
        TWITTER_SOCKET=${!i}
        ;;
        "--helm")
        ((i++))
        HELM_CMD=${!i}
        ;;
        "--kube")
        ((i++))
        KUBECTL_CMD=${!i}
        ;;
        "--mode")
        ((i++))
        MODE=${!i}
        ;;
        "-cas")
        STORAGE="cassandra"
        ;;
        "--app-key"|"-ak")
        ((i++))
        TWITTER_APP_KEY=${!i}
        ;;
        "--app-secret"|"-as")
        ((i++))
        TWITTER_APP_SECRET=${!i}
        ;;
        "--token-key"|"-tk")
        ((i++))
        TWITTER_TOKEN_KEY=${!i}
        ;;
        "--token-secret"|"-ts")
        ((i++))
        TWITTER_TOKEN_SECRET=${!i}
        ;;
        "--help"|"-h")
        echo "--kube"
        echo "--helm"
        echo "--mode"
        echo "--host"
        echo "-cas"
        echo "--app-key,-ak"
        echo "--app-secret,-as"
        echo "--token-key,-tk"
        echo "--token-secret,-ts"
        echo "--help,-h"
        ;;
    esac
done

# Edit configs
cp -r $MODE-node/* $NAMESPACE-cfg/

sed -i "s/APPKEY/$TWITTER_APP_KEY/g" $NAMESPACE-cfg/tweepub.yaml
sed -i "s/APPSECRET/$TWITTER_APP_SECRET/g" $NAMESPACE-cfg/tweepub.yaml
sed -i "s/TOKENKEY/$TWITTER_TOKEN_KEY/g" $NAMESPACE-cfg/tweepub.yaml
sed -i "s/TOKENSECRET/$TWITTER_TOKEN_SECRET/g" $NAMESPACE-cfg/tweepub.yaml

sed -i "s/KAFKARL/kafka-$NAMESPACE/g" $NAMESPACE-cfg/tweepub.yaml $NAMESPACE-cfg/tweetics.yaml
sed -i "s/ZKRELEASE/zookeeper-$NAMESPACE/g" $NAMESPACE-cfg/tweetics.yaml $NAMESPACE-cfg/kafka.yaml $NAMESPACE-cfg/spark.yaml
sed -i "s/SPARKRL/spark-$NAMESPACE/g" $NAMESPACE-cfg/tweetics.yaml

sed -i "s/STORAGE/$STORAGE/g" $NAMESPACE-cfg/tweetics.yaml $NAMESPACE-cfg/tweeviz.yaml
sed -i "s/CASRELEASE/cassandra-$NAMESPACE/g" $NAMESPACE-cfg/tweetics.yaml $NAMESPACE-cfg/tweeviz.yaml
sed -i "s/HDFSRELEASE/hdfs-$NAMESPACE/g" $NAMESPACE-cfg/tweetics.yaml $NAMESPACE-cfg/tweeviz.yaml

# Create new namespace
$KUBECTL_CMD create namespace $NAMESPACE

# Create helm charts
$HELM_CMD repo add mirantisworkloads https://mirantisworkloads.storage.googleapis.com

echo "Start charts deploying..."

$HELM_CMD install --namespace $NAMESPACE -n zookeeper-$NAMESPACE mirantisworkloads/zookeeper -f $NAMESPACE-cfg/zookeeper.yaml --wait
$HELM_CMD install --namespace $NAMESPACE -n $STORAGE-$NAMESPACE mirantisworkloads/$STORAGE -f $NAMESPACE-cfg/$STORAGE.yaml --wait

$HELM_CMD install --namespace $NAMESPACE -n kafka-$NAMESPACE mirantisworkloads/kafka -f $NAMESPACE-cfg/kafka.yaml --wait
$HELM_CMD install --namespace $NAMESPACE -n spark-$NAMESPACE mirantisworkloads/spark -f $NAMESPACE-cfg/spark.yaml --wait

$HELM_CMD install --namespace $NAMESPACE -n tweepub-$NAMESPACE mirantisworkloads/tweepub -f $NAMESPACE-cfg/tweepub.yaml --wait
$HELM_CMD install --namespace $NAMESPACE -n tweetics-$NAMESPACE mirantisworkloads/tweetics -f $NAMESPACE-cfg/tweetics.yaml --wait
$HELM_CMD install --namespace $NAMESPACE -n tweeviz-$NAMESPACE mirantisworkloads/tweeviz -f $NAMESPACE-cfg/tweeviz.yaml --wait

echo "Deploy complete!"

# Test tweeviz stats
retry_count=0

# wait a minute without check to let tweeviz collect initial data
echo "Wait a minute until tweeviz stats got initial data"
sleep 60

while [ $retry_count -le 30 ]
do
    resp=$(curl -m 10 -f $TWITTER_SOCKET/stats 2>/dev/null | jq -r '.popularity[0]')
    if [ "$resp" = null ]
    then
        echo "Tweeviz stats contain nothing, retry" $retry_count
        retry_count=$((retry_count+1))
        sleep 2
    else
        echo "Tweeviz successfully works!"
        break
    fi
done
if [ $retry_count -gt 30 ]
then
    echo "Tweeviz stats still contain nothing. Test failed."
    EXIT_CODE=1
fi

# Remove helm charts
echo "Start removing charts..."

for i in "zookeeper-$NAMESPACE" "$STORAGE-$NAMESPACE" "kafka-$NAMESPACE" "spark-$NAMESPACE" "tweepub-$NAMESPACE" "tweetics-$NAMESPACE" "tweeviz-$NAMESPACE"
do
  $HELM_CMD delete $i --purge
done

$HELM_CMD repo remove mirantisworkloads

echo "Removing charts complete."

# Remove namespace
$KUBECTL_CMD delete namespace $NAMESPACE

# Remove configs changes
rm -r $NAMESPACE-cfg/

if [ $EXIT_CODE -eq 0 ]
then
    echo "Tweeanalytics test succeed!"
fi

exit $EXIT_CODE
