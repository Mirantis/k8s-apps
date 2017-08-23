{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}

{{- define "zeppelin.fullname" -}}
{{- printf "zeppelin-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "zeppelin.spark-address" -}}
{{- if .Values.spark.deployChart -}}
    {{- $ctx := . -}}
    {{- printf "spark://" -}}
    {{- range $i, $e := until (int .Values.spark.spark.master.replicas) -}}
        {{- if $i }},{{- end -}}
        {{- printf "spark-master-%s" $.Release.Name | trunc 55 | trimSuffix "-" -}}
        {{- printf "-%d:%d" $i (int $ctx.Values.spark.spark.master.rpcPort) -}}
    {{- end -}}
{{- else -}}
    {{- printf "spark://%s" .Values.spark.externalAddress -}}
{{- end -}}
{{- end -}}

{{- define "zeppelin.address" -}}
{{ template "zeppelin.fullname" . }}:{{ .Values.port }}
{{- end -}}
