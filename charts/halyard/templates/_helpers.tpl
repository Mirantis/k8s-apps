{{- define "halyard.fullname" -}}
{{- printf "hal-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "halyard.deck.fullname" -}}
{{- printf "hal-deck-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "halyard.gate.fullname" -}}
{{- printf "hal-gate-%s" .Release.Name  | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "halyard.minio-address" -}}
    http://{{- template "minio.fullname" . -}}.{{ .Release.Namespace }}:{{ .Values.minio.port }}
{{- end -}}

{{- define "spinnaker.jenkins-address" -}}
    http://{{ template "jenkins.master-fullname" . }}:{{ .Values.jenkins.Master.port }}
{{- end -}}

{{- define "halyard.spinnaker-namespace" -}}
    {{- if .Values.prepareKubeconfig.namespace -}}
        {{- .Values.prepareKubeconfig.namespace -}}
    {{- else -}}
        {{- .Release.Namespace -}}
    {{- end -}}
{{- end -}}