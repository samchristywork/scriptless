package main

import (
	"net/http"
)

func page(w http.ResponseWriter, body string) {
	w.Write([]byte(`<!DOCTYPE html>
<head>
  <meta charset="UTF-8">
  <title>TODO</title>
  <style>
* {
  padding: 0;
  margin: 0;
}

body {
  font-family: sans;
}

button {
  margin: 0.25rem;
  padding: 12px 24px;
  font-size: 16px;
  color: #FFF;
  background-color: #07F;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  text-transform: uppercase;
  width: 100%;
}

button:hover {
  background-color: #05A;
}

input[type="text"],
input[type="password"] {
  margin: 0.25rem;
  font-size: 16px;
  border: none;
  border-radius: 0;
  outline: none;
  padding: 0.5rem;
}
  </style>
</head>
<body>
`+body+`
</body>
</html>`))
}

func table(header []string, data [][]string) string {
	content := ""
	for i:=0;i<len(data);i++ {
		content += "<tr>"
		for j:=0;j<len(data[i]);j++ {
			content += "<td>"+data[i][j]+"</td>"
		}
		content += "</tr>"
	}
	return `
<table>
	<tr>
		<th>
			`+header[0]+`
			<a href="/read?sort=name&sortdir=asc">▲</a>
			<a href="/read?sort=name&sortdir=desc">▼</a>
		</th>
		<th>
			`+header[1]+`
			<a href="/read?sort=age&sortdir=asc">▲</a>
			<a href="/read?sort=age&sortdir=desc">▼</a>
		</th>
	</tr>
	`+content+`
</table>`
}
