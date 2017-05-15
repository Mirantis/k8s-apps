{{- define "zk-fullname" -}}
{{- printf "zk-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "zookeeper-address" -}}
{{- $release := (.Release.Name | trunc 63 | trimSuffix "-") -}}
{{- range $i, $e := until (int $.Values.replicas) -}}
    {{- if $i }},{{- end -}}
         {{- printf "zk-%s-%d.zk-%s:%d" $release $i $release (int $.Values.clientPort) -}}
    {{- end -}}
{{- end -}}