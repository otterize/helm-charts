{{- define "otterize.nodeagent.fullName" -}}
otterize-node-agent
{{- end -}}
{{- define "otterize.piidetector.fullName" -}}
otterize-pii-detector
{{- end -}}
{{- define "otterize.sniffer.fullName" -}}
otterize-network-sniffer
{{- end -}}
{{- define "otterize.sniffer.securityContextConstraintsName" -}}
{{ template "otterize.sniffer.fullName" . }}-scc
{{- end -}}
{{- define "otterize.kafkawatcher.fullName" -}}
otterize-kafka-watcher
{{- end -}}
{{- define "otterize.mapper.fullName" -}}
otterize-network-mapper
{{- end -}}
{{- define "otterize.visibilitydns.fullName" -}}
otterize-visibility-dns
{{- end -}}
{{- define "otterize.mapper.configMapName" -}}
otterize-network-mapper-store
{{- end -}}
{{- define "otterize.mapper.componentConfigmap" -}}
otterize-network-mapper-component-config-map
{{- end -}}
{{ define "otterize.mapper.port" -}}
9090
{{- end -}}

{{- define "otterize.operator.apiExtraCAPath" -}}
/etc/otterize-api-extra-ca-pem
{{- end -}}

{{- define "otterize.operator.apiExtraCAPEM" -}}
{{ template "otterize.operator.apiExtraCAPath" }}/CA.pem
{{- end -}}

{{- define "otterize.networkMapper.shared_labels" -}}
app.kubernetes.io/name: network-mapper
app.kubernetes.io/part-of: otterize
app.kubernetes.io/version: {{ .Chart.Version }}
{{- with .Values.global.commonLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{- define "otterize.networkMapper.shared_pod_labels" -}}
{{- with .Values.global.podLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{- define "otterize.networkMapper.shared_annotations" -}}
app.kubernetes.io/version: {{ .Chart.Version }}
{{- with .Values.global.commonAnnotations }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{- define "otterize.networkMapper.shared_pod_annotations" -}}
{{- with .Values.global.podAnnotations }}
{{ toYaml . }}
{{- end }}
{{- end }}
