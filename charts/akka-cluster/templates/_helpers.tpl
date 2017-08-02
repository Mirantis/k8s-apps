{{- define "akka.fullname" -}}
{{- printf "akka-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "akka.producer.fullname" -}}
{{- printf "akka-producer-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "akka.consumer.fullname" -}}
{{- printf "akka-consumer-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "akka.seed.fullname" -}}
{{- printf "akka-seed-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "akka.java_opts" -}}
{{- $ctx := . -}}
{{- range $i, $e := until (int $.Values.seed.replicas) -}}
    {{- if $i }} {{ end }}
         {{- printf "-Dakka.cluster.seed-nodes.%d=akka.tcp://AkkaCluster@" $i -}}
         {{- template "akka.seed.fullname" $ctx -}}-{{ $i }}.{{- template "akka.seed.fullname" $ctx -}}:{{ $.Values.port }}
    {{- end -}}
{{- end -}}
