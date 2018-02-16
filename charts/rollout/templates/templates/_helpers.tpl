{{- define "rollout.fullname" -}}
{{- printf "roll-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "spinnaker.fullname" -}}
{{- printf "spinnaker-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "rollout.deck-address" -}}
{{ template "spinnaker.fullname" . }}-deck:{{ .Values.spinnaker.ui.port }}
{{- end -}}

{{- define "rollout.gate-address" -}}
{{ template "spinnaker.fullname" . }}-gate:8084
{{- end -}}