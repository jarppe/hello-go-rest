package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile | log.Lmicroseconds)

	host := flag.String("host", "127.0.0.1", "Host IP to listen")
	flag.Parse()

	addr := *host + ":8000"

	log.Print("Starting server at " + addr + "...")

	srv := &http.Server{
		Handler:      NewRouter(),
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
