package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func listTodos(w http.ResponseWriter, r *http.Request) {
	todos := make([]Todo, 0, len(db))
	for _, todo := range db {
		todos = append(todos, todo)
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(todos)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	todo, found := db[id]
	if found {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(todo)
	} else {
		w.WriteHeader(404)
	}
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var newTodo NewTodo
	if err := json.Unmarshal(body, &newTodo); err != nil {
		w.WriteHeader(400)
		return
	}

	id := nextID()
	todo := Todo{ID: id, Message: newTodo.Message}
	db[id] = todo

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	body, _ := ioutil.ReadAll(r.Body)
	var newTodo NewTodo
	if err := json.Unmarshal(body, &newTodo); err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode("Illegal body")
		return
	}

	todo, found := db[id]
	if found {
		todo.Message = newTodo.Message
		db[id] = todo
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(todo)
	} else {
		w.WriteHeader(404)
	}
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, found := db[id]; found {
		delete(db, id)
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
	}
}

