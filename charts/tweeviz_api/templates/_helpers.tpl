{{- define "tweeviz.api.fullname" -}}
{{- printf "tweeviz-api-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}
