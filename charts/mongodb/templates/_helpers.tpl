{{/* vim: set filetype=mustache: */}}

{{- define "cfg-name" -}}
{{ printf "mongo-cfg-%s" .Release.Name | trunc 55 | trimSuffix "-" }}
{{- end -}}

{{- define "cfg-exporter-name" -}}
{{ printf "mongo-cfg-exporter-%s" .Release.Name | trunc 55 | trimSuffix "-" }}
{{- end -}}

{{- define "shard-name" -}}
{{ printf "mongo-shard-%s" .Release.Name | trunc 55 | trimSuffix "-" }}
{{- end -}}

{{- define "shard-exporter-name" -}}
{{ printf "mongo-shard-exporter-%s" .Release.Name | trunc 55 | trimSuffix "-" }}
{{- end -}}

{{- define "router-name" -}}
{{ printf "mongo-router-%s" .Release.Name | trunc 55 | trimSuffix "-" }}
{{- end -}}

{{- define "router-exporter-name" -}}
{{ printf "mongo-router-exporter-%s" .Release.Name | trunc 55 | trimSuffix "-" }}
{{- end -}}

{{- define "mongo-name" -}}
{{ printf "mongo-%s" .Release.Name | trunc 55 | trimSuffix "-" }}
{{- end -}}

{{- define "cfg-address" -}}
{{- $ctx := . -}}
{{- range $i, $e := until (int $.Values.configServer.replicas) -}}
    {{- if $i }},{{- end -}}
    {{- template "cfg-name" $ctx -}}-{{ $i }}.{{- template "cfg-name" $ctx -}}:{{ $.Values.configServer.port }}
{{- end -}}
{{- end -}}

{{- define "router-address" -}}
{{ template "router-name" . }}:{{ .Values.router.port }}
{{- end -}}

{{- define "shard-address" -}}
{{- $ctx := . -}}
{{- range $i, $e := until (int $.Values.shard.replicas) -}}
    {{- if $i }},{{- end -}}
    {{- template "shard-name" $ctx -}}-{{ $i }}.{{- template "shard-name" $ctx -}}:{{ $.Values.shard.port }}
{{- end -}}
{{- end -}}
