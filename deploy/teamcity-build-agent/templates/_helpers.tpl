{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "teamcity-build-agent.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}-deploy
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "teamcity-build-agent.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "teamcity-build-agent.labels" -}}
helm.sh/chart: {{ include "teamcity-build-agent.chart" . }}
{{ include "teamcity-build-agent.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "teamcity-build-agent.selectorLabels" -}}
app.kubernetes.io/name: {{ include "teamcity-build-agent.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

