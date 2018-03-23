{{- define "halyard.fullname" -}}
{{- printf "hal-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "halyard.minio-address" -}}
    http://{{- template "minio.fullname" . -}}.{{ .Release.Namespace }}:{{ .Values.minio.port }}
{{- end -}}