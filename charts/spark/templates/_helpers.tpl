{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}

{{- define "master-fullname" -}}
{{- printf "spark-master-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "master-external" -}}
{{- printf "spark-master-ext-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "worker-fullname" -}}
{{- printf "spark-worker-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "zeppelin-fullname" -}}
{{- printf "zeppelin-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "spark-address" -}}
{{- $ctx := . -}}
{{- range $i, $e := until (int .Values.spark.master.replicas) -}}
    {{- if $i }},{{- end -}}
    {{- template "master-fullname" $ctx -}}
    {{- printf "-%d." $i -}}
    {{- template "master-fullname" $ctx -}}
    {{- printf ":%d" (int $ctx.Values.spark.master.rpcPort) -}}
{{- end -}}
{{- end -}}

{{- define "zookeeper-address" -}}
    {{- if .Values.zookeeper.deployChart -}}
        {{- $release := (.Release.Name | trunc 63 | trimSuffix "-") -}}
        {{- range $i, $e := until (int $.Values.zookeeper.replicas) -}}
            {{- if $i }},{{- end -}}
            {{- printf "zk-%s-%d.zk-%s:%d" $release $i $release (int $.Values.zookeeper.clientPort) -}}
        {{- end -}}
    {{- else -}}
        {{- printf .Values.zookeeper.externalAddress -}}
    {{- end -}}
{{- end -}}
