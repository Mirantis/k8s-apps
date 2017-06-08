{{- define "fullname" -}}
{{- printf "grafana-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "grafana-address" -}}
{{ template "fullname" . }}:{{ .Values.port }}
{{- end -}}

{{- define "prometheus-datasource" -}}
{{- if .Values.prometheus.deployChart -}}
{{- $url := printf "http://prometheus-server-%s:%d" (.Release.Name | trunc 63 | trimSuffix "-") (int .Values.prometheus.server.port) -}}
{{- printf `"{\"name\":\"prometheus\",\"type\":\"prometheus\",\"url\":\"%s\",\"access\":\"proxy\",\"isDefault\":true}"` $url -}}
{{- else -}}
{{- $url := (.Values.prometheus.externalAddress) -}}
{{- printf `"{\"name\":\"prometheus\",\"type\":\"prometheus\",\"url\":\"%s\",\"access\":\"proxy\",\"isDefault\":true}"` $url -}}
{{- end -}}
{{- end -}}

