{{- define "fission-fullname" -}}
{{- printf "fission-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "fission-controller" -}}
{{- printf "fission-controller-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "fission-router" -}}
{{- printf "fission-router-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "fission-poolmgr" -}}
{{- printf "fission-poolmgr-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "fission-kubewatcher" -}}
{{- printf "fission-kubewatcher-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "fission-logger" -}}
{{- printf "fission-logger-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "fission-ui" -}}
{{- printf "fission-ui-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "fission-etcd" -}}
{{- if .Values.etcd.deployChart -}}
{{- $ctx := . -}}
{{- range $i, $e := until (int $.Values.etcd.replicas) -}}
  {{- if $i }},{{- end -}}
  http://{{- template "etcd.fullname" $ctx -}}
  {{- printf "-%d." $i -}}
  {{- template "etcd.fullname" $ctx -}}
  {{- printf ":%d" (int $.Values.etcd.clientPort) -}}
{{- end -}}
{{- else -}}
{{- printf "%s" .Values.etcd.externalAddress -}}
{{- end -}}
{{- end -}}