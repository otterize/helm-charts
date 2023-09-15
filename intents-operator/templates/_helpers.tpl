{{- define "otterize.operator.tlsPath" -}}
/etc/otterize-spire
{{- end -}}

{{- define "otterize.operator.cert" -}}
{{ template "otterize.operator.tlsPath" }}/cert.pem
{{- end -}}

{{- define "otterize.operator.key" -}}
{{ template "otterize.operator.tlsPath" }}/key.pem
{{- end -}}

{{- define "otterize.operator.ca" -}}
{{ template "otterize.operator.tlsPath" }}/ca.pem
{{- end -}}

{{- define "otterize.operator.apiExtraCAPath" -}}
/etc/otterize-api-extra-ca-pem
{{- end -}}

{{- define "otterize.operator.apiExtraCAPEM" -}}
{{ template "otterize.operator.apiExtraCAPath" }}/CA.pem
{{- end -}}

{{- define "otterize.operator.mode" -}}
    {{- if not (kindIs "invalid" .Values.operator.enableEnforcement) -}}
        {{- fail "`enableEnforcement` is deprecated, please use `mode` instead. Valid values for `mode`: `defaultActive` (equivalent to `enableEnforcement`=true) and `defaultShadow` (equivalent to `enableEnforcement`=false)" -}}
    {{- end -}}
    {{- if (eq "defaultActive" .Values.operator.mode) -}}
true
    {{- else if (eq "defaultShadow" .Values.operator.mode) -}}
false
    {{- else -}}
        {{- fail (printf "Valid values for `mode`: `defaultActive` and `defaultShadow`, but you specified `%s`" .Values.operator.mode) -}}
    {{- end -}}
{{- end -}}
