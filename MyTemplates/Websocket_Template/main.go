package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	//infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	ss := NewSocketServer()

	router := chi.NewRouter()
	router.HandleFunc("/connect/{tokenID}", ss.ConnectSocketHandler)

	err := http.ListenAndServe(*addr, router)
	errorLog.Fatal(err)

}
