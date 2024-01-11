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
"true"
    {{- else if (eq "defaultShadow" .Values.operator.mode) -}}
"false"
    {{- else -}}
        {{- fail (printf "Valid values for `mode`: `defaultActive` and `defaultShadow`, but you specified `%s`" .Values.operator.mode) -}}
    {{- end -}}
{{- end -}}

{{- define "otterize.operator.allowExternalTraffic" -}}
    {{- if or (not (kindIs "invalid" .Values.operator.autoCreateNetworkPoliciesForExternalTraffic) ) (not (kindIs "invalid" .Values.operator.autoCreateNetworkPoliciesForExternalTrafficDisableIntentsRequirement) ) -}}
        {{- fail "`autoCreateNetworkPoliciesForExternalTraffic` is deprecated, please use `allowExternalTraffic` instead. \nValid values for `allowExternalTraffic`: \n\t`off` \t\t\t(equivalent to `autoCreateNetworkPoliciesForExternalTraffic`=false) \n\t`ifBlockedByOtterize` \t(equivalent to `autoCreateNetworkPoliciesForExternalTraffic`=true) \n\t`always` \t\t(equivalent to `autoCreateNetworkPoliciesForExternalTrafficDisableIntentsRequirement`=true)" -}}
    {{- end -}}
    {{- if (eq "off" .Values.operator.allowExternalTraffic) -}}
"off"
    {{- else if (eq "always" .Values.operator.allowExternalTraffic) -}}
"always"
    {{- else if (eq "ifBlockedByOtterize" .Values.operator.allowExternalTraffic) -}}
"if-blocked-by-otterize"
    {{- else -}}
        {{- fail (printf "Valid values for `allowExternalTraffic`: `off`, `ifBlockedByOtterize` and `always`, but you specified `%s`" .Values.operator.allowExternalTraffic) -}}
    {{- end -}}
{{- end -}}
