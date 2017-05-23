{{- define "fullname" -}}
{{- printf "helm-broker-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "broker-cfg" -}}
{{- printf "helm-broker-config-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

