package main

import (
	"net/http"
)

const rootPageTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
  <title>CFF Staff Info Demo</title>
  <style>
body {
  padding: 1em 5%;
  font-family: "Palatino Linotype", "Book Antiqua", Palatino, serif;
  background-color: #c0d8f0;
  text-align: center;
}

img {
  width: 85%
}

  </style>
</head>
<body>

<h1>CFF Info App Architecture</h1>
<img src="photos/app-diagram-domains.png" alt="app-topology"/>

<h2><a href="/random">Random CFF Staff Info</a></h2>

</body>
</html>
`

type rootPageHandler struct{}

func (h rootPageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rootPageTemplate))
}
