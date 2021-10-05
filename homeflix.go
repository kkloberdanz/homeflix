// Copyright (c) 2021 Kyle Kloberdanz

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

var roots []string

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func listRoot(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<pre>\n"))
	var allFiles []string
	for _, root := range roots {
		files, err := ioutil.ReadDir(root)
		if err != nil {
			fmt.Printf("error: could not list directory")
			return
		}

		for _, f := range files {
			allFiles = append(allFiles, f.Name())
		}
	}
	sort.Strings(allFiles)
	for _, fname := range allFiles {
		line := fmt.Sprintf("<a href=\"%s\">%s</a>\n", fname, fname)
		for _, root := range roots {
			path := fmt.Sprintf("%s/%s", root, fname)
			if isDirectory(path) {
				line = fmt.Sprintf("<a href=\"%s\">%s/</a>\n", fname, fname)
			}
		}
		w.Write([]byte(line))
	}
	w.Write([]byte("</pre>\n"))
	return
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	in := strings.TrimPrefix(r.URL.Path, "/")
	if strings.Contains(in, "..") {
		w.WriteHeader(http.StatusBadRequest)
		msg := "it is forbidden to include '..' in a file path on this server"
		w.Write([]byte(msg))
		return
	}

	if in == "" {
		listRoot(w)
		return
	}

	for _, root := range roots {
		filename := fmt.Sprintf("%s/%s", root, in)
		if pathExists(filename) {
			http.ServeFile(w, r, filename)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	var port = flag.String("port", "8080", "Define what TCP port to bind to")
	flag.Parse()

	roots = flag.Args()
	if len(roots) == 0 {
		roots = append(roots, ".")
	}

	address := "0.0.0.0:" + *port
	http.HandleFunc("/", handleRoot)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
