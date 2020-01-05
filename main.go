package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Todo struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type NewTodo struct {
	Message string `json:"message"`
}

type Todos map[string]Todo

func makeIdGen() func() string {
	var nextID = 0
	return func() string {
		nextID++
		return strconv.Itoa(nextID)
	}
}

var nextID = makeIdGen()
var db = Todos{}

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

func requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile | log.Lmicroseconds)

	host := flag.String("host", "127.0.0.1", "Host IP to listen")
	flag.Parse()
	addr := *host + ":8000"

	log.Print("Starting server at " + addr + "...")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", listTodos).Methods("GET")
	router.HandleFunc("/", addTodo).Methods("POST")
	router.HandleFunc("/{id}", getTodo).Methods("GET")
	router.HandleFunc("/{id}", updateTodo).Methods("POST")
	router.HandleFunc("/{id}", deleteTodo).Methods("DELETE")
	router.Use(requestLoggingMiddleware)

	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
