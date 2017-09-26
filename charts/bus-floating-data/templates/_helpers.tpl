{{- define "bfd.fullname" -}}
{{- printf "bfd-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "bfd.kafka-address" -}}
    {{- if .Values.kafka.deployChart -}}
        {{- $ctx := . -}}
        {{- range $i, $e := until (int $.Values.kafka.replicas) -}}
            {{- if $i }},{{- end -}}
            {{- template "kafka-fullname" $ctx -}}
            {{- printf "-%d." $i -}}
            {{- template "kafka-fullname" $ctx -}}
            {{- printf ":%d" (int $.Values.kafka.port) -}}
        {{- end -}}
    {{- else -}}
        {{- .Values.kafka.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "bfd.spark-address" -}}
    {{- if .Values.spark.deployChart -}}
        {{- $ctx := . -}}
        {{- range $i, $e := until (int .Values.spark.spark.master.replicas) -}}
            {{- if $i }},{{- end -}}
            {{- template "master-fullname" $ctx -}}
            {{- printf "-%d:%d" $i (int $ctx.Values.spark.spark.master.rpcPort) -}}
        {{- end -}}
    {{- else -}}
        {{- .Values.spark.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "bfd.cassandra-host" -}}
    {{- if .Values.cassandra.deployChart -}}
        {{- printf "cassandra-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
    {{- else -}}
        {{- .Values.cassandra.externalHost -}}
    {{- end -}}
{{- end -}}

{{- define "bfd.cassandra-port" -}}
    {{- if .Values.cassandra.deployChart -}}
        {{ .Values.cassandra.config.ports.cql }}
    {{- else -}}
        {{- .Values.cassandra.externalPort -}}
    {{- end -}}
{{- end -}}