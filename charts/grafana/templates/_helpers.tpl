{{- define "fullname" -}}
{{- printf "grafana-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}