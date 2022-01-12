package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

//HandleFuncを使用する場合
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

//Handleを使用する場合
type countHandler struct {
    mu sync.Mutex // guards n
    n  int
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.n++
    fmt.Fprintf(w, "count is %d\n", h.n)
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/count", new(countHandler))

	//http.ListenAndServe(":8080", nil)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
