package cg

import (
	"html/template"
)

var (
	rootTmpl *template.Template
)

const rootTmplStr = `<!DOCTYPE html>
<html>
<head>
<style>
html,body {
  margin: 0;
  padding: 10px;
  font-family: monospace;
}
table {
  border-collapse: collapse;
}
tr.even {
  background-color: #ffffff;
}
tr.odd {
  background-color: #e5e5e5;
}
td {
  padding: 8px;
  border: 1px solid #000;
}
td.blank {
  border: 0;
}
td.meven {
  background-color: #008cba;
  color: white;
}
td.modd {
  background-color: #23355c;
  color: white;
}
td.m1even {
  background-color: #005470;
  color: white;
}
td.m1odd {
  background-color: #4060a6;
  color: white;
}
td.jheading {
  border: 0;
  background-color: #000014;
  color: white;
  font-weight: bold;
  text-align: center;
}
</style>
</head>
<body>
<h2>Clebsch-Gordan Coefficients for j1 = {{ .J1 }}, j2 = {{ .J2 }}</h2>
<table>
  <tr>
    <td>m</td>
    <td>m1</td>
    <td>m2</td>
    {{- range $secIdx, $sec := .Sections }}
      {{- if $sec.PrintHeading }}
    <td class="jheading">j = {{ $sec.M }}</td>
      {{- end }}
  </tr>
      {{- $rowspan := (len $sec.Rows) }}
      {{- range $rowIdx, $row := .Rows }}
  <tr class="{{ if $secIdx | isEven }}even{{ else }}odd{{ end }}">
        {{- if eq $rowIdx 0 }}
    <td rowspan="{{ $rowspan }}" class="{{ if $secIdx | isEven }}meven{{ else }}modd{{ end }}">{{ $sec.M }}</td>
        {{- end }}
    <td class="{{ if $secIdx | isEven }}m1even{{ else }}m1odd{{ end }}">{{ $row.M1 }}</td>
    <td class="{{ if $secIdx | isEven }}m1even{{ else }}m1odd{{ end }}">{{ $row.M2 }}</td>
        {{- range $row.Values }}
    <td>{{ . }}</td>
        {{- end }}
      {{- end }}
    {{- end }}
  </tr>
</table>
</body>
</html>
`

func init() {
	funcMap := template.FuncMap{
		"isEven": func(n int) bool { return n%2 == 0 },
	}
	rootTmpl = template.Must(template.New("root").Funcs(funcMap).Parse(rootTmplStr))
}
