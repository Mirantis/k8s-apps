local ps = import "ps-lib/ps.libsonnet";
local k = import "ksonnet.beta.2/k.libsonnet";

local container = k.apps.v1beta1.statefulSet.mixin.spec.template.spec.containersType;
local psContainer = ps.platform.container;
local psStatefulSet = ps.platform.statefulSet;
local podDisruptionBudget = k.policy.v1beta1.podDisruptionBudget;
local service = k.core.v1.service;
local configMap = k.core.v1.configMap;

# common values
local namespace = "zk-kafka";
local podAntiAffinity = "hard"; # or hard or null
local persistence = "emptyDir"; # or hostPath or PersistentVolumeClaim
local hostPath = "";
local storageClass = "";
local prometheusExporter = false; # or true
local logCollector = false; # or true
local logstashHost = ["logstash-elk:5043"];

# zk
local zkReplicas = 3;
local zkFullname = "zk-ksonnet";

# kafka
local kafkaReplicas = 3;
local kafkaFullname = "kafka-ksonnet";

# templates
local zkInitContainer =
  container.new("zookeeper-init", "busybox") +
  container.command(["sh", "-c", "chown -R 1000:1000 /var/lib/zookeeper"]) +
  container.volumeMounts(container.volumeMountsType.new("datadir", "/var/lib/zookeeper"));

local zkContainer =
  psContainer.new("zookeeper", "mirantisworkloads/zookeeper:3.5.3-rc1") +
  container.ports([container.portsType.newNamed("client", 2181),
                     container.portsType.newNamed("server", 2888),
                     container.portsType.newNamed("leader-election", 3888)]) +
  (if prometheusExporter then container.ports(container.portsType.newNamed("jmx", 7071)) else {}) +
  container.env([
    container.envType.new("ZOO_LOG4J_PROP", "INFO,CONSOLE,ROLLINGFILE"),
    container.envType.new("ZK_REPLICAS", zkReplicas + ""),
    container.envType.new("ZK_LOG_LEVEL", "INFO"),
    container.envType.new("ZK_CLIENT_PORT", "2181"),
    container.envType.new("ZK_SERVER_PORT", "2888"),
    container.envType.new("ZK_ELECTION_PORT", "3888")]) +
  container.command("entrypoint.sh") +
  psContainer.readinessProbe(container.mixin.readinessProbe.exec.command("zkCheck.sh"), 15, 5) +
  psContainer.livenessProbe(container.mixin.livenessProbe.exec.command("zkCheck.sh"), 15, 5) +
  container.volumeMounts([container.volumeMountsType.new("datadir", "/var/lib/zookeeper"),
                          container.volumeMountsType.new("zoo-cfg", "/opt/zookeeper/configmap")]) +
  psContainer.logVolumeMount("/var/log/zookeeper", logCollector);

local zkStatefulSet =
  psStatefulSet.new(zkFullname, zkReplicas) +
  psStatefulSet.setMeta(zkFullname,
                        {app: zkFullname},
                        (if prometheusExporter then {"prometheus.io/scrape": "true"} else {})) +
  psStatefulSet.antiAffinity(podAntiAffinity, zkFullname) +
  psStatefulSet.initContainers(zkInitContainer, persistence != "emptyDir") +
  psStatefulSet.containers(zkContainer) +
  psStatefulSet.containers(ps.agents.logging.container.new("/var/log/zookeeper", logCollector)) +
  psStatefulSet.volumes(persistence, "datadir", hostPath, storageClass) +
  psStatefulSet.mixinVolumes({name: "zoo-cfg"} + {configMap: {name: zkFullname + "-cm"}}) +
  psStatefulSet.mixinVolumes(ps.agents.logging.volume.new(zkFullname + "-cm", logCollector));

  local zkPodDisruptionBudget =
    podDisruptionBudget.new() +
    {metadata: {name: zkFullname + "-policy"}} +
    {metadata+: {labels: {app: zkFullname}}} +
    podDisruptionBudget.mixin.spec.minAvailable(2) +
    {spec+: {selector: {matchLabels: {app: zkFullname}}}};

  local zkService =
    service.new(zkFullname,
                {app: zkFullname},
                [service.mixin.spec.portsType.port(2181) + service.mixin.spec.portsType.name("client"),
                 service.mixin.spec.portsType.port(2888) + service.mixin.spec.portsType.name("server"),
                 service.mixin.spec.portsType.port(3888) + service.mixin.spec.portsType.name("leader-election")]) +
    {spec+: {clusterIP: "None"}} +
    {metadata+: {labels: {app: zkFullname}}};

  local zkConfigMap =
    configMap.new() +
    {metadata: {name: zkFullname + "-cm"}} +
    {metadata+: {labels: {app: zkFullname}}} +
    configMap.data({"zoo.cfg":
       "tickTime=2000
        initLimit=10
        syncLimit=5
        autopurge.purgeInterval=1
        autopurge.snapRetainCount=3
        maxClientCnxns=60
        dataDir=/var/lib/zookeeper/data
        dataLogDir=/var/lib/zookeeper/log
        standaloneEnabled=false
        dynamicConfigFile=/var/lib/zookeeper/conf/zoo.cfg.dynamic
        4lw.commands.whitelist=*
        reconfigEnabled=true
        skipACL=yes",
      "java.env":
       "JVMFLAGS=\"-Xmx1G -Xms1G\"\n" +
        (if prometheusExporter then
        "JVMFLAGS=\"$JVMFLAGS -javaagent:/opt/zookeeper/jmx_prometheus_javaagent-0.9.jar=7071:/opt/zookeeper/configmap/zookeeper.yaml\"\n" else "") +
        "ZOO_LOG_DIR=\"/var/log/zookeeper\"",
      "log4j.properties":
       "zookeeper.root.logger=CONSOLE
        zookeeper.console.threshold=INFO
        zookeeper.log.maxfilesize=256MB
        zookeeper.log.maxbackupindex=20
        zookeeper.log.dir=/var/log/zookeeper
        zookeeper.log.file=zookeeper.log
        zookeeper.log.threshold=INFO
        log4j.rootLogger=${zookeeper.root.logger}
        log4j.appender.CONSOLE=org.apache.log4j.ConsoleAppender
        log4j.appender.CONSOLE.Threshold=${zookeeper.console.threshold}
        log4j.appender.CONSOLE.layout=org.apache.log4j.PatternLayout
        log4j.appender.CONSOLE.layout.ConversionPattern=%d{ISO8601} [myid:%X{myid}] - %-5p [%t:%C{1}@%L] - %m%n

        log4j.appender.ROLLINGFILE=org.apache.log4j.RollingFileAppender
        log4j.appender.ROLLINGFILE.Threshold=${zookeeper.log.threshold}
        log4j.appender.ROLLINGFILE.File=${zookeeper.log.dir}/${zookeeper.log.file}
        log4j.appender.ROLLINGFILE.MaxFileSize=${zookeeper.log.maxfilesize}
        log4j.appender.ROLLINGFILE.MaxBackupIndex=${zookeeper.log.maxbackupindex}
        log4j.appender.ROLLINGFILE.layout=org.apache.log4j.PatternLayout
        log4j.appender.ROLLINGFILE.layout.ConversionPattern=%d{ISO8601} [myid:%X{myid}] - %-5p [%t:%C{1}@%L] - %m%n"
    }) +
    (if prometheusExporter then
    configMap.data({
      "zookeeper.yaml":
       "rules:
        - pattern: \"org.apache.ZooKeeperService<name0=ReplicatedServer_id(\\d)><>(\\w+)\"
          name: \"zookeeper_$2\"
        - pattern: \"org.apache.ZooKeeperService<name0=ReplicatedServer_id(\\d), name1=replica.(\\d)><>(\\w+)\"
          name: \"zookeeper_$3\"
          labels:
            replicaId: \"$2\"
        - pattern: \"org.apache.ZooKeeperService<name0=ReplicatedServer_id(\\d), name1=replica.(\\d), name2=(\\w+)><>(\\w+)\"
          name: \"zookeeper_$4\"
          labels:
            replicaId: \"$2\"
            memberType: \"$3\"
        - pattern: \"org.apache.ZooKeeperService<name0=ReplicatedServer_id(\\d), name1=replica.(\\d), name2=(\\w+), name3=(\\w+)><>(\\w+)\"
          name: \"zookeeper_$4_$5\"
          labels:
            replicaId: \"$2\"
            memberType: \"$3\""}) else {}) +
    (if logCollector then
    configMap.data({
      "filebeat-yml":
      "filebeat.prospectors:
        - input_type: log
          paths:
            - \"/var/log/zookeeper/*.log\"
          fields:
            log_host: \"zookeeper\"
          multiline:
            pattern: '^[[:space:]]+|^Caused by:'
            negate: false
            match: after
       output.logstash:
         hosts: " + logstashHost + "
       path:
         data: \"/usr/share/filebeat/data\"
         home: \"/usr/share/filebeat\""
    }) else {});

local kafkaInitPerms =
  container.new("set-perm", "busybox") +
  container.command(["sh", "-c", "chown -R 1000:1000 /var/lib/kafka"]) +
  container.volumeMounts(container.volumeMountsType.new("datadir", "/var/lib/kafka"));

local kafkaInitWait =
  container.new("wait-for-zk", "mirantisworkloads/kafka:0.10.2.0") +
  container.command(["sh", "-c", "waitForZK.sh " + zkReplicas +" " + zkFullname]);

local kafkaContainer =
  psContainer.new("kafka", "mirantisworkloads/kafka:0.10.2.0") +
  container.ports(container.portsType.newNamed("port", 9092)) +
  (if prometheusExporter then container.ports(container.portsType.newNamed("jmx", 7071)) else {}) +
  container.env([
    container.envType.new("KAFKA_PORT", "9092"),
    container.envType.new("ZK_CONNECT", std.join(",", std.makeArray(zkReplicas, function(x) zkFullname + "-" + x + "." + zkFullname + ":2181"))),
    container.envType.new("KAFKA_HEAP_OPTS", "-Xmx1G -Xms1G")]) +
  (if prometheusExporter then container.env(container.envType.new("KAFKA_OPTS", "-javaagent:/opt/kafka/jmx_prometheus_javaagent-0.9.jar=7071:/opt/kafka/configmap/kafka-jmx.yml")) else {}) +
  container.command("entrypoint.sh") +
  psContainer.livenessProbe(container.mixin.livenessProbe.tcpSocket.port(9092), 15, 5) +
  psContainer.readinessProbe(container.mixin.readinessProbe.tcpSocket.port(9092), 15, 5) +
  container.volumeMounts([container.volumeMountsType.new("datadir", "/var/lib/kafka")]) +
  (if prometheusExporter then
    container.volumeMounts([container.volumeMountsType.new("kafka-jmx", "/opt/kafka/configmap")]) else {}) +
  psContainer.logVolumeMount("/opt/kafka/logs", logCollector);

local kafkaStatefulSet =
  psStatefulSet.new(kafkaFullname, kafkaReplicas) +
  psStatefulSet.setMeta(kafkaFullname,
                        {app: kafkaFullname},
                        (if prometheusExporter then {"prometheus.io/scrape": "true"} else {})) +
  psStatefulSet.antiAffinity(podAntiAffinity, kafkaFullname) +
  psStatefulSet.initContainers(kafkaInitPerms, persistence != "emptyDir") +
  psStatefulSet.initContainers(kafkaInitWait) +
  psStatefulSet.containers(kafkaContainer) +
  psStatefulSet.containers(ps.agents.logging.container.new("/opt/kafka/logs", logCollector)) +
  (if prometheusExporter then psStatefulSet.mixinVolumes({name: "kafka-jmx"} + {configMap: {name: kafkaFullname + "-cm"}}) else {}) +
  psStatefulSet.volumes(persistence, "datadir", hostPath, storageClass) +
  psStatefulSet.mixinVolumes(ps.agents.logging.volume.new(kafkaFullname + "-cm", logCollector));

local kafkaService =
  service.new(kafkaFullname,
              {app: kafkaFullname},
              [service.mixin.spec.portsType.port(9092) + service.mixin.spec.portsType.name("port")]) +
  {spec+: {clusterIP: "None"}} +
  {metadata+: {labels: {app: kafkaFullname}}};

local kafkaConfigMap =
  configMap.new() +
  {metadata: {name: kafkaFullname + "-cm"}} +
  {metadata+: {labels: {app: kafkaFullname}}} +
  configMap.data({"kafka-jmx.yml": "lowercaseOutputName: true
  rules:
  - pattern : kafka.cluster<type=(.+), name=(.+), topic=(.+), partition=(.+)><>Value
    name: kafka_cluster_$1_$2
    labels:
      topic: \"$3\"
      partition: \"$4\"
  - pattern : kafka.log<type=Log, name=(.+), topic=(.+), partition=(.+)><>Value
    name: kafka_log_$1
    labels:
      topic: \"$2\"
      partition: \"$3\"
  - pattern : kafka.controller<type=(.+), name=(.+)><>(Count|Value)
    name: kafka_controller_$1_$2
  - pattern : kafka.network<type=(.+), name=(.+)><>Value
    name: kafka_network_$1_$2
  - pattern : kafka.network<type=(.+), name=(.+)PerSec, request=(.+)><>Count
    name: kafka_network_$1_$2_total
    labels:
      request: \"$3\"
  - pattern : kafka.network<type=(.+), name=(\\w+), networkProcessor=(.+)><>Count
    name: kafka_network_$1_$2
    labels:
      request: \"$3\"
    type: COUNTER
  - pattern : kafka.network<type=(.+), name=(\\w+), request=(\\w+)><>Count
    name: kafka_network_$1_$2
    labels:
      request: \"$3\"
  - pattern : kafka.network<type=(.+), name=(\\w+)><>Count
    name: kafka_network_$1_$2
  - pattern : kafka.server<type=(.+), name=(.+)PerSec\\w*, topic=(.+)><>Count
    name: kafka_server_$1_$2_total
    labels:
      topic: \"$3\"
  - pattern : kafka.server<type=(.+), name=(.+)PerSec\\w*><>Count
    name: kafka_server_$1_$2_total
    type: COUNTER

  - pattern : kafka.server<type=(.+), name=(.+), clientId=(.+), topic=(.+), partition=(.*)><>(Count|Value)
    name: kafka_server_$1_$2
    labels:
      clientId: \"$3\"
      topic: \"$4\"
      partition: \"$5\"
  - pattern : kafka.server<type=(.+), name=(.+), topic=(.+), partition=(.*)><>(Count|Value)
    name: kafka_server_$1_$2
    labels:
      topic: \"$3\"
      partition: \"$4\"
  - pattern : kafka.server<type=(.+), name=(.+), topic=(.+)><>(Count|Value)
    name: kafka_server_$1_$2
    labels:
      topic: \"$3\"
    type: COUNTER

  - pattern : kafka.server<type=(.+), name=(.+), clientId=(.+), brokerHost=(.+), brokerPort=(.+)><>(Count|Value)
    name: kafka_server_$1_$2
    labels:
      clientId: \"$3\"
      broker: \"$4:$5\"
  - pattern : kafka.server<type=(.+), name=(.+), clientId=(.+)><>(Count|Value)
    name: kafka_server_$1_$2
    labels:
      clientId: \"$3\"
  - pattern : kafka.server<type=(.+), name=(.+)><>(Count|Value)
    name: kafka_server_$1_$2

  - pattern : kafka.(\\w+)<type=(.+), name=(.+)PerSec\\w*><>Count
    name: kafka_$1_$2_$3_total
  - pattern : kafka.(\\w+)<type=(.+), name=(.+)PerSec\\w*, topic=(.+)><>Count
    name: kafka_$1_$2_$3_total
    labels:
      topic: \"$4\"
    type: COUNTER
  - pattern : kafka.(\\w+)<type=(.+), name=(.+)PerSec\\w*, topic=(.+), partition=(.+)><>Count
    name: kafka_$1_$2_$3_total
    labels:
      topic: \"$4\"
      partition: \"$5\"
    type: COUNTER
  - pattern : kafka.(\\w+)<type=(.+), name=(.+)><>(Count|Value)
    name: kafka_$1_$2_$3_$4
    type: COUNTER
  - pattern : kafka.(\\w+)<type=(.+), name=(.+), (\\w+)=(.+)><>(Count|Value)
    name: kafka_$1_$2_$3_$6
    labels:
      \"$4\": \"$5\"",
  "filebeat-yml": "filebeat.prospectors:
    - input_type: log
      paths:
      - /opt/kafka/logs/*.log
      fields:
        log_host: \"kafka\"
  output.logstash:
    hosts:
    {{- range .Values.logCollector.logstashHost }}
    - {{ . | quote }}
    {{- end }}
  path:
    data: \"/usr/share/filebeat/data\"
    home: \"/usr/share/filebeat\""});

{
  zookeeper: [zkStatefulSet, zkPodDisruptionBudget, zkService, zkConfigMap],
  kafka: [kafkaStatefulSet, kafkaService, kafkaConfigMap]
}
