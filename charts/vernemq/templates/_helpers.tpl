{{- define "vernemq.fullname" -}}
{{- printf "vernemq-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

