package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strconv"
)

type Page struct {
	Title string
	Body  []byte
}

type Tasks []Task

type Task struct {
	Name, Text string
}

func (tasks *Tasks) load() error {
	filename := "tasks.json"
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, tasks)
}

func (tasks *Tasks) save() error {
	bytes, err := json.MarshalIndent(tasks, "", "    ")
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
	var task Task
	var tasks Tasks

	readJson(req, &task)

	err := tasks.load()
	if err != nil {
		panic(err)
	}
	tasks = append(tasks, task)

	tasks.save()

	bytes, err := json.Marshal(tasks)
	if err != nil {
		panic(err)
	}
	rw.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(bytes)
}

func readJson(r *http.Request, v interface{}) {
	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, v)

	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/tasks/new", newTaskHandler)
	http.HandleFunc("/tasks/", tasksHandler)
	http.ListenAndServe(":8000", nil)
}
