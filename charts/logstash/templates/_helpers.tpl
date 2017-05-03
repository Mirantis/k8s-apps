{{- define "component-fullname" -}}
{{- printf "%s-%s" .Values.component .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}


