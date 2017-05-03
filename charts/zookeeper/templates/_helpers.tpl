{{- define "zk-fullname" -}}
{{- printf "zk-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}