{{- define "fullname" -}}
{{- printf "tweeviz-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}