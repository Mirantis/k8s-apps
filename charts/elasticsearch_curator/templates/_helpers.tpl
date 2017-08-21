{{- define "curator-fullname" -}}
{{- printf "curator-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "curator-cm-fullname" -}}
{{- printf "curator-cm-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end }}

{{- define "elasticsearch-url" -}}
{{- printf "%s:%d" .Values.elasticsearch.host ( int .Values.elasticsearch.port ) }}
{{- end -}}

