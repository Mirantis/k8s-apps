{{- define "fullname" -}}
{{- printf "tweepub-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "kafka-address" -}}
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

{{- define "tweepub.hdfs-address" -}}
    {{- if .Values.hdfs.deployChart -}}
        {{- printf "hdfs-namenode-%s" .Release.Name  | trunc 63 | trimSuffix "-" }}-0.{{ printf "hdfs-namenode-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}:{{ .Values.hdfs.namenode.ui.port }}
    {{- else -}}
        {{- printf "%s" .Values.hdfs.externalAddress -}}
    {{- end -}}
{{- end -}}

