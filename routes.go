package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
	Name    string
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Path).
			Handler(route.Handler).
			Name(route.Name)
	}

	router.Use(requestLoggingMiddleware)

	return router
}

var routes = Routes{
	Route{"GET", "/", listTodos, "listTodos"},
	Route{"POST", "/", addTodo, "addTodo"},
	Route{"GET", "/{id}", getTodo, "getTodo"},
	Route{"POST", "/{id}", updateTodo, "updateTodo"},
	Route{"DELETE", "/{id}", deleteTodo, "deleteTodo"},
}

func requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}
