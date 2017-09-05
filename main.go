package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

func main() {
	http.Handle("/", http.RedirectHandler("/index.html", http.StatusMovedPermanently))
	http.Handle("/index.html", http.HandlerFunc(index))
	http.Handle("/id", http.HandlerFunc(handler))
	fmt.Println(http.ListenAndServe(":9001", nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`<html>
		<head>
			<script src="/id"></script>
		</head>
		<body>
			<h1>Test</h1>
			<pre id="cookies"></pre>
			<script>document.getElementById("cookies").innerText = document.cookie</script>
		</body>
	</html>`))
}

func handler(w http.ResponseWriter, req *http.Request) {
	// Ensure it's a GET request.
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Redirect: Seed the ID.
	if len(req.URL.Query()) == 0 {
		q := req.URL.Query()
		q.Set("id", strconv.FormatUint(rand.Uint64(), 36))

		w.Header().Set("Location", req.URL.Path+"?"+q.Encode())
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}

	// Write the ID as a script.
	w.Header().Set("Content-Type", "application/javascript")
	for h, vs := range req.URL.Query() {
		c, _ := json.Marshal(h + "=" + vs[0])
		fmt.Fprintf(w, `document.cookie = %s;`, c)
	}
}
