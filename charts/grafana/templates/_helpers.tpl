{{- define "grafana.fullname" -}}
{{- printf "grafana-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "grafana-address" -}}
{{ template "grafana.fullname" . }}:{{ .Values.port }}
{{- end -}}

{{- define "prometheus-datasource" -}}
{{- if .Values.prometheus.deployChart -}}
{{- $address := printf "prometheus-server-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- $url := printf "%s:%d" $address (int .Values.prometheus.server.port) -}}
{{- printf `"{\"name\":\"prometheus\",\"type\":\"prometheus\",\"url\":\"http://%s\",\"access\":\"proxy\",\"isDefault\":%s}"` $url (default "false" .Values.prometheus.default) -}}
{{- else -}}
{{- $url := (.Values.prometheus.externalAddress) -}}
{{- printf `"{\"name\":\"prometheus\",\"type\":\"prometheus\",\"url\":\"%s\",\"access\":\"proxy\",\"isDefault\":%s}"` $url (default "false" .Values.prometheus.default) -}}
{{- end -}}
{{- end -}}

{{- define "influxdb-datasource" -}}
{{- if .Values.influxdb.deployChart -}}
{{- $address := printf "influx-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- $url := printf "%s:%d" $address (int .Values.influxdb.ports.http.port) -}}
{{- printf `"{\"name\":\"influxdb\",\"type\":\"influxdb\",\"url\":\"http://%s\",\"access\":\"proxy\",\"database\":\"%s\",\"isDefault\":%s}"` $url (default "defaultDb" .Values.influxdb.dbInit.dbName) (default "false" .Values.influxdb.default) -}}
{{- else -}}
{{- $url := (.Values.influxdb.externalAddress) -}}
{{- printf `"{\"name\":\"influxdb\",\"type\":\"influxdb\",\"url\":\"%s\",\"access\":\"proxy\",\"database\":\"%s\",\"isDefault\":%s}"` $url .Values.influxdb.databaseName (default "false" .Values.influxdb.default ) -}}
{{- end -}}
{{- end -}}

{{- define "graphite-datasource" -}}
{{- if .Values.graphite.deployChart -}}
{{- $address := printf "graphite-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- $url := printf "%s:%d" $address (int .Values.graphite.ports.webInterface.port) -}}
{{- printf `"{\"name\":\"graphite\",\"type\":\"graphite\",\"url\":\"http://%s\",\"access\":\"proxy\",\"isDefault\":%s}"` $url (default "false" .Values.graphite.default) -}}
{{- else -}}
{{- $url := (.Values.influxdb.externalAddress) -}}
{{- printf `"{\"name\":\"graphite\",\"type\":\"graphite\",\"url\":\"http://%s\",\"access\":\"proxy\",\"isDefault\":%s}"` $url (default "false" .Values.graphite.default) -}}
{{- end -}}
{{- end -}}
