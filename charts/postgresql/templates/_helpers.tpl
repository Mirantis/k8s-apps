{{- define "postgres-fullname" -}}
{{- printf "postgresql-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "postgres-pvc-fullname" -}}
{{- printf "postgresql-pvc-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "postgres-address" -}}
{{ template "postgres-fullname" . }}:{{ .Values.port }}
{{- end -}}
