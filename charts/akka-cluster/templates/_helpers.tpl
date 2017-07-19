{{- define "akka.fullname" -}}
{{- printf "akka-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "akka.producer.fullname" -}}
{{- printf "akka-producer-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "akka.consumer.fullname" -}}
{{- printf "akka-consumer-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "akka.seed.fullname" -}}
{{- printf "akka-seed-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "akka.java_opts" -}}
{{- $release := (.Release.Name | trunc 63 | trimSuffix "-") -}}
{{- range $i, $e := until (int $.Values.seed.replicas) -}}
    {{- if $i }} {{ end }}
         {{- printf "-Dakka.cluster.seed-nodes.%d=akka.tcp://AkkaCluster@akka-seed-%s-%d.akka-seed-%s:%d" $i $release $i $release (int $.Values.port) -}}
    {{- end -}}
{{- end -}}
