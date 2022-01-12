package main

import (
	"log"
	"net/http"
	"time"
	"flag"
	"strings"
	"path"
	"path/filepath"
	"os"
	"golang.org/x/time/rate"
	"github.com/go-chi/chi/v5"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	ss :=  socketServer{
		socketMessageBuffer: 16,
		sockets: make(map[string]map[*socket]struct{}),
		socketLimiter: rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}
	//http.Handle("/", http.FileServer(http.Dir(".")))
	//http.HandleFunc("/connect/", ss.ConnectSocketHandler)
	//http.HandleFunc("/send/{roomID}", ss.SendHandler)
	//log.Fatal(http.ListenAndServe(":8000", nil))
	router := chi.NewRouter()
	FileServer(router, "/", "./ui")
	router.HandleFunc("/connect/", ss.ConnectSocketHandlerHeader)
	//router.HandleFunc("/connect/{roomID}", ss.ConnectSocketHandler)
	router.HandleFunc("/controler/connect/" ss.ConnectSocketControllerHandler)
	router.Post("/send/{socketID}", ss.SendHandler)
	err := http.ListenAndServe(*addr, router)
	log.Fatal(err)
}

func FileServer(r chi.Router, public string, static string) {

	if strings.ContainsAny(public, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	root, _ := filepath.Abs(static)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		panic("Static Documents Directory Not Found")
	}

	fs := http.StripPrefix(public, http.FileServer(http.Dir(root)))

	if public != "/" && public[len(public)-1] != '/' {
		r.Get(public, http.RedirectHandler(public+"/", 301).ServeHTTP)
		public += "/"
	}

	r.Get(public+"*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file := strings.Replace(r.RequestURI, public, "/", 1)
		if _, err := os.Stat(root + file); os.IsNotExist(err) {
			http.ServeFile(w, r, path.Join(root, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	}))
}