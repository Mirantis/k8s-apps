{{- define "fullname" -}}
{{- printf "logstash-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "cm-fullname" -}}
{{- printf "logstash-cm-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
