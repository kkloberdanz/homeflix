package todo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const form = `
<form action="/todo/post" method="post">
  <label for="todo_item">Todo Item</label><br>
  <input type="text" id="todo_item" name="todo_item"><br>
  <input type="submit" value="Submit">
</form>

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
`

const end = `
</pre>
</body>
</html>
`

func errMsg(err error, w http.ResponseWriter) {
	msg := fmt.Sprintf("error: %v", err)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
}

func Post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		errMsg(err, w)
		return
	}
	body := string(buf)
	todoItem := strings.Split(body, "=")[1]
	todoItem = strings.Replace(todoItem, "+", " ", -1)
	fmt.Printf("todo: %s\n", todoItem)
	addTodo(todoItem)
	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}

func addTodo(todoItem string) error {
	f, err := os.OpenFile("todos.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	line := fmt.Sprintf("%s\n", todoItem)
	if _, err = f.WriteString(line); err != nil {
		return err
	}

	return nil
}

func Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(form))
}
