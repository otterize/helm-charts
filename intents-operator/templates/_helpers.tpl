{{- define "otterize.operator.tlsPath" -}}
/etc/otterize-spire
{{- end -}}

{{- define "otterize.operator.cert" -}}
{{ template "otterize.operator.tlsPath" }}/svid.pem
{{- end -}}

{{- define "otterize.operator.key" -}}
{{ template "otterize.operator.tlsPath" }}/key.pem
{{- end -}}

{{- define "otterize.operator.ca" -}}
{{ template "otterize.operator.tlsPath" }}/bundle.pem
{{- end -}}