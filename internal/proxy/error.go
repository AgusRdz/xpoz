package proxy

import (
	"html/template"
	"net/http"
)

type errorData struct {
	Code    int
	Title   string
	Message string
}

var errorTmpl = template.Must(template.New("error").Parse(errorPageHTML))

func serveErrorPage(w http.ResponseWriter, code int, title, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	_ = errorTmpl.Execute(w, errorData{Code: code, Title: title, Message: message})
}

const errorPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>xpoz — {{.Code}}</title>
  <style>
    *{box-sizing:border-box;margin:0;padding:0}
    body{font-family:system-ui,sans-serif;background:#0f0f0f;color:#e5e5e5;
         display:flex;align-items:center;justify-content:center;min-height:100vh}
    .card{text-align:center;padding:2rem;max-width:480px}
    h1{font-size:5rem;font-weight:700;color:#f97316;line-height:1}
    h2{font-size:1.2rem;font-weight:500;margin:.75rem 0 1.5rem;color:#a3a3a3}
    p{line-height:1.7;color:#737373}
    footer{margin-top:2rem;font-size:.8rem;color:#404040}
  </style>
</head>
<body>
  <div class="card">
    <h1>{{.Code}}</h1>
    <h2>{{.Title}}</h2>
    <p>{{.Message}}</p>
    <footer>xpoz local tunnel</footer>
  </div>
</body>
</html>`
