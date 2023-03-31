package main

import "html/template"

const tmplStr = `<!DOCTYPE html>
<html>
<head>
<script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
<script type="text/javascript" id="MathJax-script" async
  src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-chtml.js">
</script>
</head>
<body>
$$
\begin{align}
{{ . }}
\end{align}
$$
</body>
</html>`

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("root").Parse(tmplStr))
}
