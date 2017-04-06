{{- define "fullname" -}}
{{- printf "%s-%s" .Values.component .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "cm-fullname" -}}
{{- printf "cm-%s-%s" .Values.component .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{- define "ingress-fullname" -}}
{{- printf "ingress-%s-%s" .Values.component .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}