package handler

import (
	"fmt"
	"net/http"
	//"rest-api-temp-chi/auth"
	"rest-api-temp-chi/model"
	//"time"
	//"errors"
	//"strconv"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(model.User)
	//test := r.Context().Value("test").(auth.Test)

	fmt.Printf("%v\n", user)
	fmt.Println(*user.Email)
	//fmt.Println(test)
	w.Write([]byte(fmt.Sprintf("Welcome %s!", *user.Email)))
}