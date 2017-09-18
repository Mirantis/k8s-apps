{{- define "zipkin.fullname" -}}
{{- printf "zipkin-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "zipkin.cassandra-address" -}}
{{- $ctx := . -}}
{{- range $i, $e := until (int $.Values.cassandra.config.clusterSize) -}}
  {{- if $i }},{{- end -}}
  {{- template "cassandra.fullname" $ctx -}}
  {{- printf "-%d." $i -}}
  {{- template "cassandra.fullname" $ctx -}}
  {{- printf ":%d" (int $.Values.cassandra.config.ports.cql) -}}
{{- end -}}
{{- end -}}

{{- define "zipkin.elasticsearch-address" -}}
{{ template "es-client-fullname" . }}:{{ .Values.elasticsearch.port }}
{{- end -}}