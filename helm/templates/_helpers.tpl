{{- define "kubernetes-live-images.labels" -}}
app.kubernetes.io/name: kubernetes-live-images
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}
