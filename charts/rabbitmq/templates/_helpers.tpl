{{- define "rabbitmq.fullname" -}}
{{- printf "rabbitmq-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "rabbitmq-management.fullname" -}}
{{- printf "rabbitmq-management-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "rabbitmq.address" -}}
{{- $ctx := . -}}
{{- range $i, $e := until (int $.Values.replicas) -}}
    {{- if $i }},{{- end -}}
         {{- template "rabbitmq.fullname" $ctx -}}-{{ $i }}.{{- template "rabbitmq.fullname" $ctx -}}:{{ $.Values.amqpPort }}
    {{- end -}}
{{- end -}}
