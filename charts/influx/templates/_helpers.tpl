{{- define "influx-fullname" -}}
{{- printf "influx-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "influx-cm-fullname" -}}
{{- printf "indlux-cm-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end }}
