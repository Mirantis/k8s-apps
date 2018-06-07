{{- define "es-data-fullname" -}}
{{- printf "elasticsearch-data-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "es-data-cm-fullname" -}}
{{- printf "elasticsearch-data-cm-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "es-master-fullname" -}}
{{- printf "elasticsearch-master-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "es-master-cm-fullname" -}}
{{- printf "elasticsearch-master-cm-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "es-client-fullname" -}}
{{- printf "elasticsearch-client-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "es-client-cm-fullname" -}}
{{- printf "elasticsearch-client-cm-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "es-fullname" -}}
{{- printf "elasticsearch-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "es-cm-fullname" -}}
{{- printf "elasticsearch-cm-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "es-address" -}}
{{ template "es-client-fullname" . }}:{{ .Values.port }}
{{- end -}}
