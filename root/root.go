package root

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

var Roots []string

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
<div>
    |<a href="/todo">Add Todo</a>|
    |<a href="/todos.txt">All Todos</a>|
</div>
<pre>
`))
	var allFiles []string
	for _, root := range Roots {
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
		for _, root := range Roots {
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

func Serve(w http.ResponseWriter, r *http.Request) {
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

	for _, root := range Roots {
		filename := fmt.Sprintf("%s/%s", root, in)
		if pathExists(filename) {
			http.ServeFile(w, r, filename)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
