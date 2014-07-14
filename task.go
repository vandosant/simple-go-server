package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

type Page struct {
	Title string
	Body  []byte
}

type Task struct {
	Name string
}

func (p *Page) save() error {
	filename := p.Title + ".json"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".json"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/tasks/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "%s", p.Body)
}

func newTaskHandler(w http.ResponseWriter, r *http.Request) {
	p := "Accepting new tasks here"
	fmt.Fprintf(w, p)
}

func main() {
	task1 := Task{"Get milk"}
	task2 := Task{"Get bread"}
	jsonTask1, _ := json.Marshal(task1)
	jsonTask2, _ := json.Marshal(task2)

	p1 := &Page{Title: "list1", Body: jsonTask1}
	p1.save()

	p2 := &Page{Title: "list2", Body: jsonTask2}
	p2.save()

	http.HandleFunc("/tasks/", tasksHandler)
	http.HandleFunc("/tasks/new", newTaskHandler)
	http.ListenAndServe(":8080", nil)
}
