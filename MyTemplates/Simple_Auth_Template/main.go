package main

import (
	"log"
	"net/http"
)

func main() {
	// "Signin" and "Signup" are handler that we will implement
	http.HandleFunc("/signup", SignUp)
	http.HandleFunc("/confirm", ConfirmSignUp)
	http.HandleFunc("/signin", SignIn)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}
