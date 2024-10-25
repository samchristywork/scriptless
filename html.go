package main

import (
	"net/http"
	"net/url"
	"fmt"
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

table {
	width: 100%;
	border-collapse: collapse;
	margin: 20px 0;
	font-size: 16px;
	font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	box-shadow: 0 2px 3px rgba(0,0,0,0.1);
}

th, td {
	padding: 12px 15px;
	text-align: left;
	border-bottom: 1px solid #dddddd;
}

tr:nth-child(even) {
	background-color: #f9f9f9;
}

th {
	background-color: #f1f1f1;
	font-weight: bold;
	color: #333;
}

th a {
	color: grey;
	text-decoration: none;
	font-size: 0.5rem;
	display: block;
}

tr:hover {
	background-color: #f1f1f1;
}
	</style>
</head>
<body>
` + body + `
</body>
</html>`))
}

func setPage(rawUrl string, n int) string {
	u, err := url.Parse(rawUrl)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Set("page", fmt.Sprintf("%d", n))
	u.RawQuery = q.Encode()
	return u.String()
}

func table(header []string, data [][]string, url string) string {
	maxSize := 3

	rows := min(len(data), maxSize)
	pages := float64(len(data))/float64(maxSize)

	content := ""
	for i := 0; i < rows; i++ {
		content += "<tr>"
		for j := 0; j < len(data[i]); j++ {
			content += "<td>" + data[i][j] + "</td>"
		}
		content += "</tr>"
	}

	links := ""
	for i := 0; float64(i) < pages; i++ {
		links += fmt.Sprintf(`<a href="%s">%d</a>`,
		setPage(url, i+1), i+1)
	}

	return `
<table>
	<tr>
		<th>
			<span>` + header[0] + `</span>
			<span style="display: inline-block;">
				<a href="/read?sort=name&sortdir=asc">▲</a>
				<a href="/read?sort=name&sortdir=desc">▼</a>
			</span>
		</th>
		<th>
			<span>` + header[1] + `</span>
			<span style="display: inline-block;">
				<a href="/read?sort=age&sortdir=asc">▲</a>
				<a href="/read?sort=age&sortdir=desc">▼</a>
			</span>
		</th>
	</tr>
	` + content + `
</table>` + links
}

func div(style string, body string) string {
	return `<div style="` + style + `">` + body + `</div>`
}

func centeredBox(body string) string {
	return div(`
	padding: 1rem;
	background-color: #eee;
	border-radius: 0.25rem;
	position: absolute;
	top: 50vh;
	left: 50vw;
	transform: translate(-50%, -50%);
	box-shadow: 0 4px 8px rgba(0, 0, 0, 0.4);
	text-align: center;
	`, body)
}
