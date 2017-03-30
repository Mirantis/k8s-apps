{{/*
Create fully qualified names. We truncate at 24 chars
because some Kubernetes name fields are limited to
this (by the DNS naming spec).
*/}}
{{- define "data-fullname" -}}
{{- printf "%s-data-%s" .Release.Name .Values.Elasticsearch.Name | trunc 24 -}}
{{- end -}}


{{- define "master-fullname" -}}
{{- printf "%s-master-%s" .Release.Name .Values.Elasticsearch.Name | trunc 24 -}}
{{- end -}}

{{- define "client-fullname" -}}
{{- printf "%s-client-%s" .Release.Name .Values.Elasticsearch.Name | trunc 24 -}}
{{- end -}}

{{- define "es-fullname" -}}
{{- printf "%s-%s" .Release.Name .Values.Elasticsearch.Name | trunc 24 -}}
{{- end -}}

