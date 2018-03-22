{{- define "halyard.fullname" -}}
{{- printf "hal-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}
