package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)


func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", home)
	router.HandleFunc("/snippet", showSnippet)
	router.HandleFunc("/articles/{category}", articlesCategoryHandler)
	router.HandleFunc("/articles/{category}/{id:[0-9]+}", articleHandler)
	//router.HandleFunc("/snippet/create", createSnippet)
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", router)
	log.Fatal(err)
}

