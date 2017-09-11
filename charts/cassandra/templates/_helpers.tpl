{{- define "cassandra.fullname" -}}
{{- printf "cassandra-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "cassandra.address" -}}
{{ template "cassandra.fullname" . }}:{{ .Values.config.ports.cql }}
{{- end -}}
