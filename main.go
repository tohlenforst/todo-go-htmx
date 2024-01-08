package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"
)

type Todo struct {
	Id        int
	Text      string
	Completed bool
}

const todo_html = `
	<div style="border: 1px solid black;" hx-target="closest div" hx-swap="outerHTML">
		<div>ID: {{.Id}}</div>
		<div>Text: {{.Text}}</div>
		<div class="status">Completed?: {{.Completed}}</div>
		<button hx-put="/api/todos?id={{.Id}}">Toggle Completed</button>
		<button hx-delete="/api/todos?id={{.Id}}">Delete</button>
	</div>
`

func main() {
	todos := []Todo{
		{1, "Learn Go", false},
		{2, "Learn Vue", false},
		{3, "Learn React", false},
	}

	var todo_template = template.Must(template.New("todo").Parse(todo_html))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/api/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			for _, todo := range todos {
				todo_template.Execute(w, todo)
			}
		}
		if r.Method == http.MethodPost {
			newTodo := Todo{todos[len(todos)-1].Id + 1, r.FormValue("newTodo"), false}
			todos = append(todos, newTodo)
			todo_template.Execute(w, newTodo)
		}
		if r.Method == http.MethodPut {
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				http.Error(w, "", http.StatusBadRequest)
			}
			oldTodo := slices.IndexFunc(todos, func(todo Todo) bool {
				return todo.Id == id
			})
			todos[oldTodo].Completed = !todos[oldTodo].Completed
			todo_template.Execute(w, todos[oldTodo])
		}
		if r.Method == http.MethodDelete {
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				http.Error(w, "", http.StatusBadRequest)
			}
			todos = slices.DeleteFunc(todos, func(todo Todo) bool {
				return todo.Id == id
			})
		}
	})

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
