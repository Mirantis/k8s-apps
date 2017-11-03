{{- define "gerrit.fullname" -}}
{{- printf "gerrit-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}
