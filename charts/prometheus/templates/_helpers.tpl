{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a fully qualified alertmanager name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "alertmanager.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "prometheus-alertmanager-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a fully qualified kube-state-metrics name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "kubeStateMetrics.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "prometheus-kube-state-metrics-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a fully qualified node-exporter name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "nodeExporter.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "prometheus-node-exporter-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a fully qualified Prometheus server name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "server.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "prometheus-server-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "alertmanager.address" -}}
{{ template "alertmanager.fullname" . }}:{{ .Values.alertmanager.port }}
{{- end -}}

{{- define "kubeStateMetrics.address" -}}
{{ template "kubeStateMetrics.fullname" . }}:{{ .Values.kubeStateMetrics.port }}
{{- end -}}

{{- define "nodeExporter.address" -}}
{{ template "nodeExporter.fullname" . }}:{{ .Values.nodeExporter.port }}
{{- end -}}

{{- define "server.address" -}}
{{ template "server.fullname" . }}:{{ .Values.server.port }}
{{- end -}}

