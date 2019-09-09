package main

import (
	"bytes"
	"html/template"
	"io"

	"code.cloudfoundry.org/lager"
)

const pageTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
  <style>
body {
  padding: 1em 10%; 
  font-family: "Palatino Linotype", "Book Antiqua", Palatino, serif;
}

h1 {
	font-size: 3.5em;
	margin: 0.2em 0 0.4em;
}

h2 {
	font-size: 2.5em;
	margin: 0.5em 0;
}

p {
	font-size: smaller;
}

.metadata {
	font-family: "Lucida Console", Monaco, monospace;
}

.metadata tr {
	padding: 0.5em 0;
}

.metadata td:first-child {
	font-weight: bold;
	padding: 1px 1em 1px 0;
}

.headshot {
	float: left;
	width: 40%;
	padding: 0;
	border: 2px solid gray;
	margin: 0 2em;
}
  </style>
</head>
<body>
  <img class="headshot" src="/photos/{{.Member.ID}}.png" alt="{{.Member.Name}}"/>
  <h1>{{.Member.Name}}</h1>
  <h2>{{.Member.Title}}</h2>
  <p>{{.Member.Bio}}</p>

<div class="metadata">
	<h3>Request Metadata</h3>

	<table>
	 <tr><td>Duration</td><td>{{.Metadata.Duration}}</td></tr>
	 <tr><td>Member URL</td><td>{{.Metadata.URL}}</td></tr>
	 <tr><td>Connected Address</td><td>{{.Metadata.Addr}}</td></tr>
	 <tr><td>Server IP</td><td>{{.Metadata.ServerIP}}</td></tr>
	</table> 
</div>

</body>
</html>
`

type PagePresenter struct {
	template *template.Template
}

func NewPagePresenter() *PagePresenter {
	return &PagePresenter{
		template: template.Must(template.New("page").Parse(pageTemplate)),
	}
}

func (p *PagePresenter) WritePage(logger lager.Logger, w io.Writer, data FetchResult) error {
	buf := &bytes.Buffer{}

	err := p.template.Execute(buf, data)
	if err != nil {
		logger.Error("failed-to-render-template", err)
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		logger.Error("failed-to-write-rendered-template", err)
		return err
	}

	return nil
}
