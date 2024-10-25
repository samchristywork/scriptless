package main

import (
	"net/http"
	"net/url"
	"fmt"
)

func page(w http.ResponseWriter, style string, body string) {
	w.Write([]byte(`<!DOCTYPE html>
<head>
	<meta charset="UTF-8">
	<title>TODO</title>
	<style>`+style+`</style>
</head>
<body>
` + body + `
</body>
</html>`))
}

func setParam(rawUrl string, param string, value string) string {
	u, err := url.Parse(rawUrl)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Set(param, value)
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
		setParam(url, "page", fmt.Sprintf("%d", i+1)), i+1)
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
