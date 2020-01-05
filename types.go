package main

type Todo struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type NewTodo struct {
	Message string `json:"message"`
}

type Todos map[string]Todo

