{{- define "otterize.sniffer.fullName" -}}
otterize-network-sniffer
{{- end -}}
{{- define "otterize.kafkawatcher.fullName" -}}
otterize-kafka-watcher
{{- end -}}
{{- define "otterize.mapper.fullName" -}}
otterize-network-mapper
{{- end -}}
{{- define "otterize.iamlive.fullName" -}}
otterize-iamlive
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