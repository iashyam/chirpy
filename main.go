package main

//let's start

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type message struct {
	text         string
	content_type string
}

func (m message) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", m.content_type)
	w.WriteHeader(200)
	fmt.Fprint(w, m.text)
}

func (cfg *apiConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	fmt.Fprintf(w, "<html>\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n<p>Chirpy has been visited %v times!</p>\n</body>\n</html>", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) middleWareInc(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		handle.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middleWareReset(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Store(0)
		handle.ServeHTTP(w, r)
	})
}

const port = "8080"

func main() {

	cfg := apiConfig{}
	m := message{text: "OK", content_type: "text/plain; charset=utf-8"}

	fmt.Println("Server listening at ")
	serverMux := http.NewServeMux()
	filepathRoot := ""
	// serverMux.HandleFunc("/healthz", m.ServeHttp)
	serverMux.Handle("GET /api/healthz", m)
	serverMux.Handle("GET /admin/metrics", &cfg)
	serverMux.Handle("/app/", cfg.middleWareInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	serverMux.Handle("POST /admin/reset", cfg.middleWareReset(m))
	serverMux.Handle("/app/assets/logo.png", cfg.middleWareInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	localServer := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}
	err := localServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
