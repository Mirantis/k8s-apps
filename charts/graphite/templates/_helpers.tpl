{{- define "graphite-fullname" -}}
{{- printf "graphite-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "graphite-cm-fullname" -}}
{{- printf "graphite-cm-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end }}
