{{- define "kafka-fullname" -}}
{{- printf "kafka-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "zookeeper-fullname" -}}
{{- printf "zk-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "zk-address" -}}
    {{- if .Values.zookeeper.deployChart -}}
        {{- $release := (.Release.Name | trunc 63 | trimSuffix "-") -}}
        {{- range $i, $e := until (int $.Values.zookeeper.replicas) -}}
            {{- if $i }},{{- end -}}
            {{- printf "zk-%s-%d.zk-%s:%d" $release $i $release (int $.Values.zookeeper.clientPort) -}}
        {{- end -}}
    {{- else -}}
        {{- printf "%s" .Values.zookeeper.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "kafka.address" -}}
{{ $replicas := int (.Values.replicas)}}
{{- range $i, $e := until $replicas -}}
    {{- if $i -}},{{- end -}}
    kafka-{{ $.Release.Name | trunc 63 | trimSuffix "-" }}-{{ $i }}.kafka-{{ $.Release.Name | trunc 63 | trimSuffix "-" }}:{{ $.Values.port }}
{{- end -}}
{{- end -}}
