{{- define "postgres-fullname" -}}
{{- printf "postgresql-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "postgres-pvc-fullname" -}}
{{- printf "postgresql-pvc-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}