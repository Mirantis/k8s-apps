{{- define "vernemq.fullname" -}}
{{- printf "vernemq-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

