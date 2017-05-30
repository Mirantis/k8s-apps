{{- define "fullname" -}}
{{- printf "kafka-broker-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

