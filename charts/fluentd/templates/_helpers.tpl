{{- define "fluent-fullname" -}}
{{- printf "fluentd-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "fluent-cm-fullname" -}}
{{- printf "fluentd-cm-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end }}

{{- define "elasticsearch-url" -}}
{{- printf "%s:%d" .Values.elasticsearch.host ( int .Values.elasticsearch.port ) }}
{{- end -}}

