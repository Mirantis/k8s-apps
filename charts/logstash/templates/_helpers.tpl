{{- define "ls-fullname" -}}
{{- printf "logstash-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "ls-cm-fullname" -}}
{{- printf "logstash-cm-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "ls-es-url" -}}
    {{- if .Values.elasticsearch.external -}}
        {{- printf "%s:%d" .Values.elasticsearch.host ( int .Values.elasticsearch.port ) }}
    {{- else -}}
        {{- printf "%s:%d" ( include "es-fullname" . ) ( int .Values.elasticsearch.port ) }}
    {{- end -}}
{{- end -}}

{{- define "ls-address" -}}
{{ template "ls-fullname" . }}:{{ .Values.port }}
{{- end -}}
