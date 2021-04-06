
{{- if and .Values.barometerAgent.apiKey .Values.barometerAgent.apiKeySecret.secretName -}}
    {{- required "Only one of barometerAgent.apiKey and barometerAgent.apiKeySecret may be supplied!" nil -}}
{{- end -}}

{{- if and (not .Values.barometerAgent.apiKey) (not .Values.barometerAgent.apiKeySecret.secretName) -}}
    {{- required "One of barometerAgent.apiKey and barometerAgent.apiKeySecret must be supplied!" nil -}}
{{- end -}}