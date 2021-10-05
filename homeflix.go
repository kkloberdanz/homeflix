// httpserver.go
package main

import (
	"flag"
	"net/http"
)

func main() {
	var port = flag.String("port", "8080", "Define what TCP port to bind to")
	var root = flag.String("root", ".", "Define the root filesystem path")

	flag.Parse()
	panic(http.ListenAndServe("0.0.0.0:"+*port, http.FileServer(http.Dir(*root))))
}
