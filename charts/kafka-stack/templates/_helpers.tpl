{{- define "kafka-connector.fullname" -}}
{{- printf "kafka-connector-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "kafka-connector.address" -}}
{{- template "kafka-connector.fullname" . -}}:{{- printf "%d" (int $.Values.kafkaConnector.port) -}}
{{- end -}}

{{- define "kafka-rest.fullname" -}}
{{- printf "kafka-rest-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "kafka-rest.address" -}}
{{- template "kafka-rest.fullname" . -}}:{{- printf "%d" ( int .Values.kafkaRest.port ) -}}
{{- end -}}

{{- define "zk-fullname" -}}
{{- printf "zk-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "schema-registry.fullname" -}}
{{- printf "schema-registry-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "schema-registry.address" -}}
{{- template "schema-registry.fullname" . -}}:{{- printf "%d" (int $.Values.schemaRegistry.port) -}}
{{- end -}}

{{- define "kafka-stack.kafka-address" -}}
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
        {{- printf "%s" .Values.kafka.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "kafka-stack.zk-address" -}}
    {{- if .Values.kafka.deployChart -}}
        {{- if .Values.kafka.zookeeper.deployChart -}}
            {{- $ctx := . -}}
            {{- range $i, $e := until (int $.Values.kafka.zookeeper.replicas) -}}
                {{- if $i }},{{- end -}}
                {{- template "zk-fullname" $ctx -}}
                {{- printf "-%d." $i -}}
                {{- template "zk-fullname" $ctx -}}
                {{- printf ":%d" (int $.Values.kafka.zookeeper.clientPort) -}}
             {{- end -}}
        {{- else -}}
            {{- printf "%s" .Values.kafka.zookeeper.externalAddress -}}
        {{- end -}}
     {{- else -}}
         {{- printf "%s" .Values.kafka.zookeeper.externalAddress -}}
     {{- end -}}
{{- end -}}
