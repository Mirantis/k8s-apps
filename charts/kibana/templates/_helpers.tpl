{{- define "kib-fullname" -}}
{{- printf "kibana-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "kib-cm-fullname" -}}
{{- printf "kibana-cm-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{- define "kib-ingress-fullname" -}}
{{- printf "kibana-ingress-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{- define "kib-es-url" -}}
    {{- if .Values.elasticsearch.external -}}
        {{- printf "%s:%d" .Values.elasticsearch.host ( int .Values.elasticsearch.port ) }}
    {{- else -}}
        {{- printf "%s:%d" ( include "es-fullname" . ) ( int .Values.elasticsearch.port ) }}
    {{- end -}}
{{- end -}}