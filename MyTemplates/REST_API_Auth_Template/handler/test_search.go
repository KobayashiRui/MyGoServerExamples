package handler

import (
	"fmt"
	"net/http"
	//"rest-api-temp-chi/auth"
	//"rest-api-temp-chi/model"
	"encoding/json"
	//"time"
	//"errors"
	//"strconv"
)

func TestSearch(w http.ResponseWriter, r *http.Request) {
	//test := r.Context().Value("test").(auth.Test)

	decoder := json.NewDecoder(r.Body)
	fmt.Println(decoder)

	//fmt.Printf("%v\n", user)
	//fmt.Println(*user.Email)
	//fmt.Println(test)
	w.Write([]byte(fmt.Sprintf("search")))
}