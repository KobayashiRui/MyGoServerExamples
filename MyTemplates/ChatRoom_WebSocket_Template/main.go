package main

import (
	"log"
	"net/http"
	"time"
	"golang.org/x/time/rate"
)

func main() {
	ss :=  socketServer{
		socketMessageBuffer: 16,
		sockets: make(map[string]map[*socket]struct{}),
		socketLimiter: rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/connect/", ss.connectSocketHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}