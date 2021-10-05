// httpserver.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var root = flag.String("root", ".", "Define the root filesystem path")

func sendFile(filename string, w http.ResponseWriter) error {
	// open input file
	fi, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fi.Close()

	// make a read buffer
	r := bufio.NewReader(fi)

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
	}

	//if err = w.Flush(); err != nil {
	//	return err
	//}

	return nil
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	in := strings.TrimPrefix(r.URL.Path, "/")
	if strings.Contains(in, "..") {
		w.WriteHeader(http.StatusBadRequest)
		//w.Write("it is forbidden to include '..' in a file path on this server")
		return
	}
	filename := fmt.Sprintf("%s/%s", *root, in)
	fmt.Printf("serving: '%s'\n", filename)
	http.ServeFile(w, r, filename)
	//err := sendFile(filename, w)
	//if err != nil {
	//	w.WriteHeader(http.FileNotFound)
	//	fmt.Printf("not a file: '%s'\n", filename)
	//}
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
