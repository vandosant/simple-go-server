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

type Tasks []Task

type Task struct {
	Name, Text string
}

func (tasks *Tasks) load(bytes []byte) error {
	return json.Unmarshal(bytes, tasks)
}

func (tasks *Tasks) save() error {
	bytes, err := json.Marshal(tasks)
	filename := "tasks.json"
	if err != nil {
		return err
	} else {
		return ioutil.WriteFile(filename, bytes, 0600)
	}
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

func tasksHandler(w http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/tasks/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "%s", p.Body)
}

func newTaskHandler(rw http.ResponseWriter, req *http.Request) {
	var jsonStream = []byte(`[
		{"Name": "Buy Milk", "Text": "Two percent, organic free-range."},
		{"Name": "Buy Bread", "Text": "Whole wheat."},
		{"Name": "Get Camping Supplies", "Text": "Bug spray, sunscreen."},
		{"Name": "Pay Phone Bill", "Text": "Due by Saturday."},
		{"Name": "Call Mom", "Text": "You missed her birthday."}
		]`)
	var tasks Tasks
	err := tasks.load(jsonStream)
	if err != nil {
		panic(err)
	}

	tasks.save()

	bytes, err := json.Marshal(tasks)
	if err != nil {
		panic(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(bytes)
}

func main() {
	http.HandleFunc("/tasks/new", newTaskHandler)
	http.HandleFunc("/tasks/", tasksHandler)
	http.ListenAndServe(":8000", nil)
}
