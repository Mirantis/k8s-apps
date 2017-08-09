{{- define "autoscale.fullname" -}}
{{- printf "autoscale-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "autoscale.prometheus-address" -}}
    {{- if .Values.prometheus.deployChart -}}
        http://{{- template "server.fullname" . -}}:{{- .Values.prometheus.server.port -}}
    {{- else -}}
        {{- printf "%s" .Values.prometheus.externalAddress -}}
    {{- end -}}
{{- end -}}
