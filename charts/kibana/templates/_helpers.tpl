{{- define "fullname" -}}
{{- printf "kibana-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "cm-fullname" -}}
{{- printf "kibana-cm-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{- define "ingress-fullname" -}}
{{- printf "kibana-ingress-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}