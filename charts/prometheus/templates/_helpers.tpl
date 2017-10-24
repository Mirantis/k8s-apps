{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a fully qualified alertmanager name.
*/}}
{{- define "alertmanager.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "prometheus-alertmanager-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a fully qualified kube-state-metrics name.
*/}}
{{- define "kubeStateMetrics.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "prometheus-kube-state-metrics-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a fully qualified node-exporter name.
*/}}
{{- define "nodeExporter.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "prometheus-node-exporter-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a fully qualified Prometheus server name.
*/}}
{{- define "server.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "prometheus-server-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "alertmanager.address" -}}
{{- if .Values.alertmanager.deploy }}
{{ template "alertmanager.fullname" . }}.{{ .Release.Namespace }}:{{ .Values.alertmanager.port }}
{{- else }}
{{ .Values.alertmanager.externalAddress }}
{{- end }}
{{- end -}}

{{- define "kubeStateMetrics.address" -}}
{{ template "kubeStateMetrics.fullname" . }}.{{ .Release.Namespace }}:{{ .Values.kubeStateMetrics.port }}
{{- end -}}

{{- define "nodeExporter.address" -}}
{{ template "nodeExporter.fullname" . }}.{{ .Release.Namespace }}:{{ .Values.nodeExporter.port }}
{{- end -}}

{{- define "server.address" -}}
{{ template "server.fullname" . }}.{{ .Release.Namespace }}:{{ .Values.server.port }}
{{- end -}}
