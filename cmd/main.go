package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

type Todo struct {
	Id      int
	Message string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	data := map[int]Todo{
		1: {Id: 1, Message: "Buy Book"},
	}

	todosHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		templ := template.Must(template.ParseFiles("templates/index.html"))
		templ.Execute(w, data)
	}

	addTodoHandler := func(w http.ResponseWriter, r *http.Request) {
		message := r.PostFormValue("message")
		if message == "" {
			return
		}
		templ := template.Must(template.ParseFiles("templates/index.html"))
		id := getNextID(data)
		todo := Todo{Id: id, Message: message}
		data[id] = todo
		templ.ExecuteTemplate(w, "todo-list-element", todo)
	}

	deleteTodoHandler := func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return
		}
		delete(data, id)
	}

	editTodoHandler := func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return
		}
		templ := template.Must(template.ParseFiles("templates/components/edit.html"))
		templ.Execute(w, data[id])
	}

	saveTodoHandler := func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PostFormValue("id"))
		if err != nil {
			return
		}
		message := r.PostFormValue("message")
		todo := Todo{Id: id, Message: message}
		data[id] = todo
		templ := template.Must(template.ParseFiles("templates/index.html"))
		templ.ExecuteTemplate(w, "todo-list-element", todo)
	}

	http.HandleFunc("/", todosHandler)
	http.HandleFunc("/add-todo", addTodoHandler)
	http.HandleFunc("/delete-todo/{id}", deleteTodoHandler)
	http.HandleFunc("/edit-todo/{id}", editTodoHandler)
	http.HandleFunc("/save-todo", saveTodoHandler)

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func getNextID(todos map[int]Todo) int {
	maxID := 0
	for id := range todos {
		if id > maxID {
			maxID = id
		}
	}
	return maxID + 1
}
