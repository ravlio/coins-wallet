{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "service.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "service.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "nc" -}}
{{- printf "%s-%s" .name .id -}}
{{- end -}}

{{- define "cheader" }}
 image: "{{ .image.repository }}:{{ .image.tag }}"
 imagePullPolicy: {{ .image.pullPolicy }}
 resources:
  requests:
   memory: {{ .resources.requests.memory }}
   cpu: {{ .resources.requests.cpu }}
  limits:
   memory: {{ .resources.limits.memory }}
   cpu: {{ .resources.limits.cpu }}
{{- end }}

{{- define "mscnm" -}}
{{- $name := .name | replace "-" "--" -}}
{{- $id := .id | replace "-" "--" -}}
{{- printf "%s-%s" $name $id -}}
{{- end -}}

{{- define "mscnms" -}}
{{- $name := .name | replace "-" "--" -}}
{{- printf "%s" $name -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "service.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}
