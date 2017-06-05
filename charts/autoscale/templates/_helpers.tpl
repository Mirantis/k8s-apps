{{- define "autoscale.fullname" -}}
{{- printf "autoscale-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "autoscale.prometheus-address" -}}
    {{- if .Values.prometheus.deployChart -}}
        http://{{- printf "prometheus-server-%s:%d" .Release.Name ( int .Values.prometheus.server.port) | trunc 63 | trimSuffix "-" -}}
    {{- else -}}
        {{- printf "%s" .Values.prometheus.externalAddress -}}
    {{- end -}}
{{- end -}}
