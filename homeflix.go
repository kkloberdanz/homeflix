// httpserver.go
package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	in := r.URL.Path
	in = strings.TrimPrefix(in, "/")
	out := fmt.Sprintf("hello %s\n", in)
	w.Write([]byte(out))
}

func main() {
	var port = flag.String("port", "8080", "Define what TCP port to bind to")
	//var root = flag.String("root", ".", "Define the root filesystem path")

	flag.Parse()

	address := "0.0.0.0:" + *port
	http.HandleFunc("/", sayHello)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
