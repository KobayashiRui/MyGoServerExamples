package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strings"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()       //urlが渡すオプションを解析します。POSTに対してはレスポンスパケットのボディを解析します（request body）
    //注意：もしParseFormメソッドがコールされなければ、以下でフォームのデータを取得することができません。
    fmt.Println(r.Form) //これらのデータはサーバのプリント情報に出力されます
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    fmt.Fprintf(w, "Hello astaxie!") //ここでwに書き込まれたものがクライアントに出力されます。
}

func login(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) //リクエストを取得するメソッド
    if r.Method == "GET" {
        t, _ := template.ParseFiles("login.gtpl")
        t.Execute(w, nil)
    } else {
        //ログインデータがリクエストされ、ログインのロジック判断が実行されます。
		r.ParseForm()
        fmt.Println("username:", r.Form["username"])
        fmt.Println("password:", r.Form["password"])
    }
}

func main() {
    http.HandleFunc("/", sayhelloName)       //アクセスのルーティングを設定します
    http.HandleFunc("/login", login)         //アクセスのルーティングを設定します
    err := http.ListenAndServe(":9090", nil) //監視するポートを設定します
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}