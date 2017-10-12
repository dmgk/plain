package handler

import (
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(homeTmpl))
}

var homeTmpl = `
<!DOCTYPE html>
<html>
	<head>
		<link rel="shortcut icon" href="/favicon.ico" />
		<title>plain.im - a plain pastebin</title>
		<style>
			body {
				font-family: monospace;
			}
			h1 {
				font-size: 1.2em;
				font-weight: bold;
				padding: 0.7em 0 0.7em 0;
				margin: 0;
			}
			h2 {
				font-size: inherit;
				font-weight: bold;
				padding: 1.2em 0 1.2em 0;
				margin: 0;
			}
			p {
				padding: 0 3em 0 3em;
				display: block;
				margin: 0;
			}
			a, a:visited {
				color: inherit;
				text-decoration: underline;
			}
		</style>
	</head>
	<body>
		<h1>plain.im</h1>
		plain.im is a plaintext paste service with 24 hours expiration time.
		<h2>API</h2>
		<p>
			Paste the command output:<br/>
			<b>command | curl -sF "plain=<-" https://plain.im/</b><br/><br/>

			Create an alias:<br/>
			<b>alias plain='curl -sF "plain=<-" https://plain.im/'</b><br/>
			to shorten paste command to just<br/>
			<b>command | plain</b><br/><br/>

			Paste and make link available as an X selection for pasting into X applications (requires xclip):<br/>
			<b>alias plain='curl -sF "plain=<-" https://plain.im/ | tee /dev/stderr | xclip'</b>
		</p>
		<h2>PASTE REMOVAL</h2>
		<p>
			Paste can be removed by making a DELETE request to the paste URL:<br/>
			<b>curl -X DELETE https://plain.im/key'</b>
		</p>
		<h2>SOURCE CODE</h2>
		<p>
			<a href="https://github.com/dmgk/plain">https://github.com/dmgk/plain</a>
		</p>
		<h2>LICENSE</h2>
		<p>
			MIT
		</p>
	</body>
</html>
`
