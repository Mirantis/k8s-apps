{{- define "fullname" -}}
{{- printf "jenkins-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}