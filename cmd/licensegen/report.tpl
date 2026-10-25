{{- range . }}
## {{ .Name }} {{ .Version }}

- **License:** {{ .LicenseName }}
{{ if and .LicenseURL (ne .LicenseURL "Unknown") -}}
- **Source:** {{ .LicenseURL }}
{{ end }}

{{ if and .LicensePath (ne .LicensePath "Unknown") -}}
{{ markdownBlock .LicenseText }}
{{- else -}}
License text unavailable.
{{- end }}

---

{{ end -}}
