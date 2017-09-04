{{- define "etcd.fullname" -}}
{{- printf "etcd-%s" .Release.Name | trunc 55 | trimSuffix "-" -}}
{{- end -}}

{{- define "etcd.address" -}}
{{- $ctx := . -}}
{{- range $i, $e := until (int $.Values.replicas) -}}
  {{- if $i }},{{- end -}}
  {{- template "etcd.fullname" $ctx -}}
  {{- printf "-%d." $i -}}
  {{- template "etcd.fullname" $ctx -}}
  {{- printf ":%d" (int $.Values.clientPort) -}}
{{- end -}}
{{- end -}}