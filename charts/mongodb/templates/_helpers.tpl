{{/* vim: set filetype=mustache: */}}

{{- define "cfg-name" -}}
{{ printf "mongo-cfg-%s" .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end -}}
