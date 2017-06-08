{{- define "fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" $name .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "cassandra-address" -}}
{{ template "fullname" . }}:{{ .Values.config.ports.cql }}
{{- end -}}