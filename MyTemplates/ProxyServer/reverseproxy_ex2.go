package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	origin, _ := url.Parse("http://localhost:8188/")

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}

	proxy := &httputil.ReverseProxy{Director: director}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		log.Printf(token)
		if token != "hogehoge" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":9001", nil))
}
