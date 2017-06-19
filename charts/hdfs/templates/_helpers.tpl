{{- define "namenode-fullname" -}}
{{- printf "hdfs-namenode-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "datanode-fullname" -}}
{{- printf "hdfs-datanode-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "configmap-fullname" -}}
{{- printf "hdfs-configs-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "hdfs-ui-fullname" -}}
{{- printf "hdfs-ui-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "namenode-address" -}}
{{- printf "hdfs-namenode-%s" .Release.Name  | trunc 63 | trimSuffix "-" }}-0.{{ printf "hdfs-namenode-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}:{{ .Values.namenode.port }}
{{- end -}}

{{- define "hdfs-ui-address" -}}
{{ template "hdfs-ui-fullname" . }}:{{ .Values.namenode.ui.port }}
{{- end -}}
