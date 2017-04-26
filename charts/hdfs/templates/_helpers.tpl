{{- define "namenode-fullname" -}}
{{- printf "hdfs-namenode-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "datanode-fullname" -}}
{{- printf "hdfs-datanode-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "configmap-fullname" -}}
{{- printf "hdfs-configs-%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "namenode-address" -}}
{{- printf "hdfs-namenode-%s" .Release.Name  | trunc 63 | trimSuffix "-" }}-0.{{ printf "hdfs-namenode-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}
