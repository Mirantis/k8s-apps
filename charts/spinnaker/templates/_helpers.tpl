{{- define "fullname" -}}
{{- printf "spinnaker-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "spinnaker.minio-address" -}}
    {{- if .Values.minio.deployChart -}}
        http://{{- template "minio.fullname" . -}}:{{ .Values.minio.port }}
    {{- else -}}
        {{- printf "%s" .Values.minio.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "spinnaker.jenkins-master-address" -}}
    {{- if .Values.jenkins.deployChart -}}
        {{ template "jenkins.master-fullname" . }}:{{ .Values.jenkins.Master.port }}
    {{- else -}}
        {{- printf "%s" .Values.jenkins.Master.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "spinnaker.jenkins-agent-address" -}}
    {{- if .Values.jenkins.deployChart -}}
        {{ template "jenkins.agent-fullname" . }}:{{ .Values.jenkins.Agent.port }}
    {{- else -}}
        {{- printf "%s" .Values.jenkins.Agent.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "spinnaker.redis-address" -}}
    {{- if .Values.redis.deployChart -}}
        redis://{{ template "redis-fullname" . }}-0.{{ template "redis-fullname" . }}:{{ .Values.redis.config.redisPort }}
    {{- else -}}
        {{- printf "%s" .Values.redis.externalAddress -}}
    {{- end -}}
{{- end -}}
