{{- define "monocular.fullname" -}}
{{- printf "monocular-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "monocular.mongodb-address" -}}
    {{- if .Values.mongodb.deployChart -}}
        {{- template "router-name" . -}}:{{ .Values.mongodb.router.port }}
    {{- else -}}
        {{- printf "%s" .Values.mongodb.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "monocular.redis-address" -}}
    {{- if .Values.redis.deployChart -}}
        {{ template "redis-fullname" . }}-0.{{ template "redis-fullname" . }}:{{ .Values.redis.config.redisPort }}
    {{- else -}}
        {{- printf "%s" .Values.redis.externalAddress -}}
    {{- end -}}
{{- end -}}