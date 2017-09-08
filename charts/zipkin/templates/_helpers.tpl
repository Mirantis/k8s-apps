{{- define "zipkin.fullname" -}}
{{- printf "zipkin-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}
