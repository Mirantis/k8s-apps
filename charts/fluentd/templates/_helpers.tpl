{{- define "fluent-fullname" -}}
{{- printf "fluentd-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "fluent-cm-fullname" -}}
{{- printf "fluentd-cm-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end }}

{{- define "elasticsearch-url" -}}
    {{- if .Values.elasticsearch.deployChart -}}
        {{- printf "%s:%d" ( include "es-fullname" . ) ( int .Values.elasticsearch.port | default 9200 ) }}
    {{- else -}}
        {{ .Values.elasticsearch.externalAddress }}
    {{- end -}}
{{- end -}}
