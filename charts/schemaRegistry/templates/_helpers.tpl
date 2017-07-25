{{- define "schema-registry.fullname" -}}
{{- printf "schema-registry-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "schema-registry.address" -}}
{{- printf "schema-registry-%s:%d" (.Release.Name | trunc 63 | trimSuffix "-") (int $.Values.port) -}}
{{- end -}}

{{- define "schema-registry.kafka-address" -}}
    {{- if .Values.kafka.deployChart -}}
        {{- $release := (.Release.Name | trunc 63 | trimSuffix "-") -}}
        {{- range $i, $e := until (int $.Values.kafka.replicas) -}}
            {{- if $i }},{{- end -}}
            {{- printf "kafka-%s-%d.kafka-%s:%d" $release $i $release (int $.Values.kafka.port) -}}
        {{- end -}}
    {{- else -}}
        {{- printf "%s" .Values.kafka.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "schema-registry.zk-address" -}}
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
