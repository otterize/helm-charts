{{- define "otterize.operator.apiExtraCAPath" -}}
/etc/otterize-api-extra-ca-pem
{{- end -}}

{{- define "otterize.operator.apiExtraCAPEM" -}}
{{ template "otterize.operator.apiExtraCAPath" }}/CA.pem
{{- end -}}

{{- define "otterize.credentialsOperator.shared_labels" -}}
app.kubernetes.io/name: credentials-operator
app.kubernetes.io/part-of: otterize
app.kubernetes.io/version: {{ .Chart.Version }}
{{- with .Values.global.commonLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{- define "otterize.credentialsOperator.shared_pod_labels" -}}
{{- with .Values.global.podLabels }}
{{ toYaml . }}
{{- end }}
{{ if eq true .Values.global.azure.enabled }}
azure.workload.identity/use: "true"
{{ end }}
{{- end }}

{{- define "otterize.credentialsOperator.shared_annotations" -}}
app.kubernetes.io/version: {{ .Chart.Version }}
{{- with .Values.global.commonAnnotations }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{- define "otterize.credentialsOperator.shared_pod_annotations" -}}
{{- with .Values.global.podAnnotations }}
{{ toYaml . }}
{{- end }}
{{- end }}