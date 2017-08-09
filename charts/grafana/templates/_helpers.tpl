{{- define "fullname" -}}
{{- printf "grafana-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "grafana-address" -}}
{{ template "fullname" . }}:{{ .Values.port }}
{{- end -}}

{{- define "prometheus-datasource" -}}
{{- if .Values.prometheus.deployChart -}}
{{- $address := printf "prometheus-server-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- $url := printf "%s:%d" $address (int .Values.prometheus.server.port) -}}
{{- printf `"{\"name\":\"prometheus\",\"type\":\"prometheus\",\"url\":\"http://%s\",\"access\":\"proxy\",\"isDefault\":true}"` $url -}}
{{- else -}}
{{- $url := (.Values.prometheus.externalAddress) -}}
{{- printf `"{\"name\":\"prometheus\",\"type\":\"prometheus\",\"url\":\"%s\",\"access\":\"proxy\",\"isDefault\":true}"` $url -}}
{{- end -}}
{{- end -}}

