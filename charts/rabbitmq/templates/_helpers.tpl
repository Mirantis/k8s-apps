{{- define "rabbitmq.fullname" -}}
{{- printf "rabbitmq-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "rabbitmq.address" -}}
{{- $release := (.Release.Name | trunc 63 | trimSuffix "-") -}}
{{- range $i, $e := until (int $.Values.replicas) -}}
    {{- if $i }},{{- end -}}
         {{- printf "rabbitmq-%s-%d.rabbitmq-%s:%d" $release $i $release (int $.Values.amqpPort) -}}
    {{- end -}}
{{- end -}}