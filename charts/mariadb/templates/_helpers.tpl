{{- define "mariadb-fullname" -}}
{{- printf "mariadb-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "mariadb-fullname-pvc" -}}
{{- printf "mariadb-pvc-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "mariadb-address" -}}
{{ template "mariadb-fullname" . }}:{{ .Values.port }}
{{- end -}}
