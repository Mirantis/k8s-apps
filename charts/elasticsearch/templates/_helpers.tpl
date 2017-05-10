{{- define "data-fullname" -}}
{{- printf "elasticsearch-data-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}


{{- define "master-fullname" -}}
{{- printf "elasticsearch-master-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "client-fullname" -}}
{{- printf "elasticsearch-client-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "fullname" -}}
{{- printf "elasticsearch-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
