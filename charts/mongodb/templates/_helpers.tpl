{{/* vim: set filetype=mustache: */}}

{{- define "cfg-name" -}}
{{ printf "mongo-cfg-%s" .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{- define "shard-name" -}}
{{ printf "mongo-shard-%s" .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{- define "router-name" -}}
{{ printf "mongo-router-%s" .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{- define "cfg-address" -}}
{{- $release := (.Release.Name | trunc 63 | trimSuffix "-") -}}
{{- range $i, $e := until (int $.Values.configServer.replicas) -}}
    {{- if $i }},{{- end -}}
    {{- printf "mongo-cfg-%s-%d.mongo-cfg-%s:%d" $release $i $release (int $.Values.configServer.port) -}}
{{- end -}}
{{- end -}}

{{- define "router-address" -}}
{{- $release := (.Release.Name | trunc 63 | trimSuffix "-") -}}
{{ printf "mongo-router-%s:%d" $release (int .Values.router.port) }}
{{- end -}}
