{{- define "jenkins.master-fullname" -}}
{{- printf "jenkins-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "jenkins.agent-fullname" -}}
{{- printf "jenkins-agent-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "jenkins.master-address" -}}
{{ template "jenkins.master-fullname" . }}:{{ .Values.Master.port }}
{{- end -}}

{{- define "jenkins.agent-address" -}}
{{ template "jenkins.agent-fullname" . }}:{{ .Values.Agent.port }}
{{- end -}}