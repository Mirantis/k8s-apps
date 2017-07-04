{{- define "clair.fullname" -}}
{{- printf "clair-%s" .Release.Name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "clair.postgresql-connection-string" -}}
    {{- if .Values.postgresql.deployChart -}}
        host={{ printf "postgresql-%s" .Release.Name | trunc 63 | trimSuffix "-" }} port={{ .Values.postgresql.port }} user={{ .Values.postgresql.credentials.user }} password={{ .Values.postgresql.credentials.password }} sslmode=disable statement_timeout=60000
    {{- else -}}
        {{- printf "%s" .Values.postgresql.externalAddress -}}
    {{- end -}}
{{- end -}}

{{- define "clair.config" -}}
    clair:
      database:
        type: pgsql
        options:
          source: {{ template "clair.postgresql-connection-string" . }}
{{ toYaml .Values.database | indent 10 }}
      api:
{{ toYaml .Values.api | indent 8 }}
      updater:
{{ toYaml .Values.updater | indent 8 }}
      notifier:
{{ toYaml .Values.notifier | indent 8 }}
{{- end -}}