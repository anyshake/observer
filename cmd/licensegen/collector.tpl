[
{{- range $index, $library := . }}
{{- if $index }},{{ end }}
{
  "Name": {{ printf "%q" .Name }},
  "Version": {{ printf "%q" .Version }},
  "LicenseName": {{ printf "%q" .LicenseName }},
  "LicensePath": {{ printf "%q" .LicensePath }},
  "LicenseURL": {{ printf "%q" .LicenseURL }}
}
{{- end }}
]
