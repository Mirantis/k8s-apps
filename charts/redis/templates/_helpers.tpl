{{- define "redis-fullname" -}}
{{- printf "redis-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "sentinel-fullname" -}}
{{- printf "sentinel-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "sentinel-svc-fullname" -}}
{{- printf "redis-sentinel-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "cluster-svc-fullname" -}}
{{- printf "redis-cluster-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}