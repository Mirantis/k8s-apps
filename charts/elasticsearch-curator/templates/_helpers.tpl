{{- define "curator-fullname" -}}
{{- printf "curator-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "curator-cm-fullname" -}}
{{- printf "curator-cm-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end }}

{{- define "elasticsearch-fqdn" -}}
    {{- if .Values.elasticsearch.deployChart -}}
        {{ include "es-client-fullname" . }}
    {{- else -}}
        {{ .Values.elasticsearch.externalAddress }}
    {{- end -}}
{{- end -}}
