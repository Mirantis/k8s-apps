{{- define "redis-fullname" -}}
{{- printf "redis-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "sentinel-fullname" -}}
{{- printf "sentinel-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "sentinel-svc-fullname" -}}
{{- printf "redis-sentinel-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "cluster-svc-fullname" -}}
{{- printf "redis-cluster-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "redis-address" -}}
{{- $ctx := . -}}
{{- range $i, $e := until (int .Values.replicas) -}}
    {{- if $i -}},{{- end -}}
    {{- template "redis-fullname" $ctx -}}-{{ $i }}.{{- template "sentinel-svc-fullname" $ctx -}}:{{ $.Values.config.redisPort }}
{{- end -}}
{{- end -}}
