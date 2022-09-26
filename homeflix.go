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
	if err != nil {
		return false
	}

	return info.IsDir()
}

func listRoot(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport"
     content="width=device-width, initial-scale=1, user-scalable=yes">

  <title>Homeflix</title>
  <link rel="shortcut icon"
    href="http://homeflix.local/favicon.svg">
</head>

<style type="text/css" media="screen">

body {
    font-family: sans-serif;
    padding: 10px;
    background: #aaaaaa;
}

.card {
    background-color: white;
    margin-top: 1px;
    padding: 1%;
}

</style>

<body>
<pre>
`))
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
	var last string
	for _, fname := range allFiles {
		if fname == last {
			continue
		}
		line := fmt.Sprintf("<div class=\"card\"> <a href=\"%s\">%s</a> </div>\n", fname, fname)
		for _, root := range roots {
			path := fmt.Sprintf("%s/%s", root, fname)
			if isDirectory(path) {
				line = fmt.Sprintf("<div class=\"card\"> <a href=\"%s\">%s/</a> </div>\n", fname, fname)
			}
		}
		w.Write([]byte(line))
		last = fname
	}
	w.Write([]byte(`
</pre>
</body>
</html>
`))
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

	fmt.Fprintf(os.Stderr, "serving from %s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
