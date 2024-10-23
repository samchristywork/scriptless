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
