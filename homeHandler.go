package main

import "net/http"

func getHomeHandler(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusOK)
	res.Header().Add("Content-Type", "text/html")
	res.Write([]byte(`
		<html>
		<title>Cats API</title>
		<link rel="icon" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='0.9em' font-size='80'>😺</text></svg>">
		<style>
		html, body {
			width: 100%;
		}
		a {
			text-decoration: none;
		}
		</style>
		<body>
			<h2>Software version: ` + version + `</h2>
			<br/>
			<a href="swagger/"><h3>🖥️ Swagger OpenAPI UI</h3></a>
		<body>
		</html>
	`))
}
