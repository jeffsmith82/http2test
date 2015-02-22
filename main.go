// Go and Server Sent Events for HTTP/1.1 and HTTP/2.0
//go:generate go run $GOROOT/src/crypto/tls/generate_cert.go -host localhost

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bradfitz/http2"
)

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/sse", sseHandler)
	s := http.Server{Addr: ":8080"}
	http2.ConfigureServer(&s, nil)
	err := s.ListenAndServeTLS("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}
}

var html = `<html>
<head>
<script>
var ev = new EventSource("/sse");
ev.onmessage = function(event) {
	document.getElementById("ev").innerHTML = event.data;
}
</script>
</head>
<body>
<p id="ev"></p>
</body>
</html>`

func mainHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, html)
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	conn, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Oops", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.WriteHeader(http.StatusOK)
	i := 0
	for {
		_, err := fmt.Fprintf(w, "data: hello world %d\n\n", i)
		if err != nil {
			break
		}
		i++
		conn.Flush()
		time.Sleep(time.Second)
	}
}
