{{- define "es-data-fullname" -}}
{{- printf "elasticsearch-data-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}


{{- define "es-master-fullname" -}}
{{- printf "elasticsearch-master-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "es-client-fullname" -}}
{{- printf "elasticsearch-client-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "es-fullname" -}}
{{- printf "elasticsearch-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "es-address" -}}
{{ template "es-client-fullname" . }}:{{ .Values.port }}
{{- end -}}
