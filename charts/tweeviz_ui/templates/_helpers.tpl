{{- define "tweeviz.ui.fullname" -}}
{{- printf "tweeviz-ui-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}