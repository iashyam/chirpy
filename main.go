package main

//let's start

import (
	"log"
	"net/http"
)

const port="8080"

func main() {

	serverMux := http.NewServeMux()
	serverMux.Handle("/", http.FileServer(http.Dir(".")))

	localServer := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}
	err := localServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
