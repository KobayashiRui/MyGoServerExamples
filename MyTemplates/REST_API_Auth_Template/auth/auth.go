package auth

import (
	"context"
	"fmt"
	"net/http"
	"encoding/json"
	//"github.com/go-chi/chi/v5"
	"log"
	"rest-api-temp-chi/database/mongodb"
	"rest-api-temp-chi/database/redisdb"
	"rest-api-temp-chi/model"
	"time"
	//"errors"
	//"strconv"
	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"crypto/rand"
)

type AuthHandler struct {
	ErrorLog *log.Logger
	InfoLog *log.Logger 
	TempUser *mongodb.TempUserModel
	User *mongodb.UserModel
	Session *redisdb.SessionModel
}

type Test struct {
	Hoge string
	Fuga int
}

func (ah * AuthHandler) HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func (ah * AuthHandler) CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}


func (ah * AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {

	var tempUser model.TempUser
	err := json.NewDecoder(r.Body).Decode(&tempUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_temporaryToken, _ := uuid.NewRandom() 
	temporaryToken := _temporaryToken.String()
	tempUser.Token = &temporaryToken
	hashedPW, err := ah.HashPassword(*tempUser.Password)
	if(err != nil){
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tempUser.Password = &hashedPW


	insertTempUser, err := ah.TempUser.Insert(tempUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(insertTempUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//temporaryUsers[temporaryToken] = json_user
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func (ah * AuthHandler) ConfirmSignUp(w http.ResponseWriter, r *http.Request){
	//v := r.URL.Query()
	temporaryToken := r.URL.Query().Get("token")
	tempUserID := chi.URLParam(r, "userID")
	result, err := ah.TempUser.GetOne(tempUserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("confirm")
	fmt.Println(tempUserID)
	fmt.Println(*result.Token)
	fmt.Println(temporaryToken)
	if(*result.Token == temporaryToken) {
		var user model.User
		user.SetFromTempUser(result)
		ah.InfoLog.Println(user)
		ah.TempUser.DeleteOne(tempUserID)
		ah.User.Insert(user)
		w.Write([]byte("confirmed"))
		return
	}else{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (ah * AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	//var json_user jsonUser
	//Userのログイン
	//userID := chi.URLParam(r, "userID")
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(*user.Email)
	userData, err := ah.User.GetOneFromEmail(*user.Email)
	pwCheck := ah.CheckPasswordHash(*user.Password, *userData.Password)
	if !pwCheck {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//expectedPassword, ok := users[json_user.Email]
	//if !ok || expectedPassword != json_user.Password {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

	_sessionToken, _ := uuid.NewRandom()
	sessionToken := _sessionToken.String()
	//fmt.Println(userData.ID.String())


	ah.Session.SetUserID(sessionToken, userData.ID.Hex(), (120 * time.Second) )
	//caches[sessionToken] = json_user.Email
	//cookieの設定
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(110 * time.Second), //tokenの存続時間
		//Secure: true,
		HttpOnly: true,
	})
}

func (ah * AuthHandler) GetTempUsers(w http.ResponseWriter, r *http.Request) {
 	tempUserList, err := ah.TempUser.Get()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	jsonCompany, err := json.Marshal(tempUserList)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCompany)
	return
}

func (ah * AuthHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
 	userList, err := ah.User.Get()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	jsonCompany, err := json.Marshal(userList)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonCompany)
	return

}

//middleware
func (ah * AuthHandler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				fmt.Println("no token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionToken := c.Value
		fmt.Println("get auth")
		fmt.Println(sessionToken)

		userID, err := ah.Session.GetUserID(sessionToken)
		if err != nil{
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Println(userID)
		resultUser, err := ah.User.GetOne(userID)
		if err != nil{
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

	  	ctx := context.WithValue(r.Context(), "user", *resultUser)

	  	next.ServeHTTP(w, r.WithContext(ctx))
	})
}