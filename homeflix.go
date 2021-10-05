package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
)

var root = flag.String("root", ".", "Define the root filesystem path")

func handleRoot(w http.ResponseWriter, r *http.Request) {
	in := strings.TrimPrefix(r.URL.Path, "/")
	if strings.Contains(in, "..") {
		w.WriteHeader(http.StatusBadRequest)
		msg := "it is forbidden to include '..' in a file path on this server"
		w.Write([]byte(msg))
		return
	}
	filename := fmt.Sprintf("%s/%s", *root, in)
	fmt.Printf("serving: '%s'\n", filename)
	http.ServeFile(w, r, filename)
}

func main() {
	var port = flag.String("port", "8080", "Define what TCP port to bind to")

	flag.Parse()

	address := "0.0.0.0:" + *port
	http.HandleFunc("/", handleRoot)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
