package main
import (
	"fmt"
	"html/template"
	"net/http"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"time"
	//"errors"
	//"strconv"
)

type Article struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	TimeStamp time.Time `json:"time_stamp"`
}

func newArticle() *Article {
	_article := new(Article)
	_article.ID = 1
	_article.Title = ""
	_article.Content = ""
	_article.TimeStamp = time.Now()

	return _article
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	ts, err := template.ParseFiles("./ui/html/home.page.tmpl")
	if err != nil {
		app.errorLog.Println(err.Error())
		//http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
}

func (app *application) getArticle(w http.ResponseWriter, r *http.Request){
	articleID := chi.URLParam(r, "articleID")
	app.infoLog.Println(articleID)
	article_data := Article{
		ID: 1,
		Title: "hoge",
		Content: "fuga fuga content",
	}
	article_json, err := json.Marshal(article_data)
	if err != nil {
		app.serverError(w, err) // Use the serverError() helper
		return
	}
	app.infoLog.Println(string(article_json))
	w.Header().Set("Content-Type", "application/json")
	w.Write(article_json)
}

func (app *application) createArticle(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	http.Error(w, "Method Not Allowed", 405)
	//	return
	//}
	//w.Write([]byte("Create a new Articles..."))
	new_article := newArticle()
	//var unmarshalErr *json.UnmarshalTypeError	
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(new_article)
	fmt.Printf("%+v\n", *new_article) // with Variable name

	if err != nil {
		app.serverError(w, err)
		return
	}

	//if err != nil {
	//	if errors.As(err, &unmarshalErr) {
	//		//errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
	//	} else {
	//		//errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
	//	}
	//	return
	//}


	resp := make(map[string]string)
	resp["message"] = "ok"

	resp_json, err := json.Marshal(resp)
	if err != nil {
		app.serverError(w, err)
	}
	app.infoLog.Println(string(resp_json))
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp_json)

}

