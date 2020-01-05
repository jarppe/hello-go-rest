package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)


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
