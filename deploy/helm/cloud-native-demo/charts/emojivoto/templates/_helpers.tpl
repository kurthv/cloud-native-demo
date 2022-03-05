{{/*
Expand the name of the chart.
*/}}
{{- define "emojivoto.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "emojivoto.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "emojivoto.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "emojivoto.labels" -}}
helm.sh/chart: {{ include "emojivoto.chart" . }}
{{ include "emojivoto.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "emojivoto.selectorLabels" -}}
app.kubernetes.io/name: {{ include "emojivoto.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the names of the service accounts to use
*/}}

{{- define "emojivoto.serviceAccountNameEmoji" -}}
{{- if .Values.serviceAccount.create }}
{{- .Values.serviceAccount.emojiName }}
{{- else }}
{{- default "default" .Values.serviceAccount.emojiName }}
{{- end }}
{{- end }}

{{- define "emojivoto.serviceAccountNameVote" -}}
{{- if .Values.serviceAccount.create }}
{{- .Values.serviceAccount.voteName }}
{{- else }}
{{- default "default" .Values.serviceAccount.voteName }}
{{- end }}
{{- end }}

{{- define "emojivoto.serviceAccountNameWeb" -}}
{{- if .Values.serviceAccount.create }}
{{- .Values.serviceAccount.webName }}
{{- else }}
{{- default "default" .Values.serviceAccount.webName }}
{{- end }}
{{- end }}
