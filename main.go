// Copyright (c) 2021 Kyle Kloberdanz

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/kkloberdanz/homeflix/root"
	"github.com/kkloberdanz/homeflix/todo"
)

func main() {
	var port = flag.String("port", "8080", "Define what TCP port to bind to")
	flag.Parse()

	root.Roots = flag.Args()
	if len(root.Roots) == 0 {
		root.Roots = append(root.Roots, ".")
	}

	address := "0.0.0.0:" + *port
	http.HandleFunc("/", root.Serve)
	http.HandleFunc("/todo", todo.Serve)
	http.HandleFunc("/todo/post", todo.Post)

	fmt.Fprintf(os.Stderr, "serving from %s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
