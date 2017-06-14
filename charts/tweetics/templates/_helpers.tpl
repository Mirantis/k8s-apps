{{- define "tweetics.fullname" -}}
{{- printf "tweetics-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "tweetics.kafka-address" -}}
    {{- if .Values.kafka.deployChart -}}
        {{- $release := (.Release.Name | trunc 63 | trimSuffix "-") -}}
        {{- range $i, $e := until (int $.Values.kafka.replicas) -}}
            {{- if $i }},{{- end -}}
            {{- printf "kafka-%s-%d.kafka-%s:%d" $release $i $release (int $.Values.kafka.port) -}}
        {{- end -}}
    {{- else -}}
        {{- printf "%s" .Values.kafka.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "tweetics.zk-address" -}}
    {{- if .Values.zookeeper.deployChart -}}
        {{- $release := (.Release.Name | trunc 63 | trimSuffix "-") -}}
        {{- range $i, $e := until (int $.Values.zookeeper.replicas) -}}
            {{- if $i }},{{- end -}}
            {{- printf "zk-%s-%d.zk-%s:%d" $release $i $release (int $.Values.zookeeper.clientPort) -}}
        {{- end -}}
    {{- else -}}
        {{- printf "%s" .Values.zookeeper.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "tweetics.spark-address" -}}
    {{- if .Values.spark.deployChart -}}
        {{- $release := (.Release.Name | trunc 63 | trimSuffix "-") -}}
        {{- range $i, $e := until (int $.Values.spark.spark.master.replicas) -}}
            {{- if $i }},{{- end -}}
            {{- printf "spark-master-%s-%d.spark-master-%s:%d" $release $i $release (int $.Values.spark.spark.master.rpcPort) -}}
        {{- end -}}
    {{- else -}}
        {{- printf "%s" .Values.spark.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "tweetics.hdfs-address" -}}
    {{- if .Values.hdfs.deployChart -}}
        {{- printf "hdfs-namenode-%s" .Release.Name  | trunc 63 | trimSuffix "-" }}-0.{{ printf "hdfs-namenode-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}:{{ .Values.hdfs.namenode.port }}
    {{- else -}}
        {{- printf "%s" .Values.hdfs.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "tweetics.cassandra-address" -}}
    {{- if .Values.cassandra.deployChart -}}
        {{- printf "cassandra-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
    {{- else -}}
        {{- .Values.cassandra.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "tweetics-storage" -}}
    {{- if eq .Values.storage "hdfs" -}}
        {{- printf "hdfs://" -}}{{ template "tweetics.hdfs-address" . }}{{- .Values.hdfs.path -}}
    {{- else -}}
        {{ template "tweetics.cassandra-address" . }}:{{- .Values.cassandra.keyspace -}}:{{- .Values.cassandra.table -}}
    {{- end -}}
{{- end -}}
