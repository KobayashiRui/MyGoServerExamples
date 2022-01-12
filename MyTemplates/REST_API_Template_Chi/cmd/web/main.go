package main

import (
	"flag"
	"net/http"
	"log"
	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger 
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{ 
		errorLog: errorLog, 
		infoLog: infoLog,
	}
	router := chi.NewRouter()

	router.Get("/", app.home)
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	//get article list
	//router.Get("/articles", )
	//get article data
	router.Get("/articles/{articleID}", app.getArticle)
	//create article
	router.Post("/articles", app.createArticle)


	err := http.ListenAndServe(*addr, router)
	errorLog.Fatal(err)
}