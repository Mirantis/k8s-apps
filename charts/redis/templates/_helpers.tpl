{{- define "redis-fullname" -}}
{{- printf "redis-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "redis-address" -}}
{{- $ctx := . -}}
{{- range $i, $e := until (int .Values.replicas) -}}
    {{- if $i -}},{{- end -}}
    {{- template "redis-fullname" $ctx -}}-{{ $i }}.{{- template "redis-fullname" $ctx -}}:{{ $.Values.config.redisPort }}
{{- end -}}
{{- end -}}
