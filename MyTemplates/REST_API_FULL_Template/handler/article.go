package handler
import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"rest-api-temp-chi/dataservice/mongodb"
	//"rest-api-temp-chi/model"
	//"time"
	//"errors"
	//"strconv"
)

type ArticleHandler struct {
	ErrorLog *log.Logger
	InfoLog *log.Logger 
	Article *mongodb.ArticleModel
}

func (ah *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request){

	//DB

	articles, err_db := ah.Article.Get()
	fmt.Println(articles)
	if err_db != nil {
		ah.ErrorLog.Println(err_db.Error())
	}

	//articleID := chi.URLParam(r, "articleID")
	//ah.InfoLog.Println(articleID)
	//article_data := model.Article{
	//	Title: "hoge",
	//	Content: "fuga fuga content",
	//}
	article_json, err := json.Marshal(articles)
	if err != nil {
		serverError(w, err) // Use the serverError() helper
		return
	}
	ah.InfoLog.Println(string(article_json))
	w.Header().Set("Content-Type", "application/json")
	w.Write(article_json)
}

func (ah *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request){

	articleID := chi.URLParam(r, "articleID")

	//DB
	article, err_db := ah.Article.GetOne(articleID)
	//fmt.Println(*article)
	if err_db != nil {
		ah.ErrorLog.Println(err_db.Error())
		//serverError(w, err_db)
		return
	}

	//articleID := chi.URLParam(r, "articleID")
	//ah.InfoLog.Println(articleID)
	//article_data := model.Article{
	//	Title: "hoge",
	//	Content: "fuga fuga content",
	//}
	article_json, err := json.Marshal(article)
	if err != nil {
		serverError(w, err) // Use the serverError() helper
		return
	}
	ah.InfoLog.Println(string(article_json))
	w.Header().Set("Content-Type", "application/json")
	w.Write(article_json)
}

func (ah *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	w.Header().Set("Allow", http.MethodPost)
	//	http.Error(w, "Method Not Allowed", 405)
	//	return
	//}
	//w.Write([]byte("Create a new Articles..."))

	type JsonArticle struct {
		Title string `json:"title"`
		Content string `json:"content"`
	}

	//new_article := struct {
	//	title string `json:"title"`
	//	content string `json:"content"`
	//}
	//var unmarshalErr *json.UnmarshalTypeError
	new_article := new(JsonArticle)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(new_article)
	fmt.Printf("%+v\n", *new_article) // with Variable name

	if err != nil {
		serverError(w, err)
		return
	}

	ah.Article.Insert(new_article.Title, new_article.Content)

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
		serverError(w, err)
	}
	ah.InfoLog.Println(string(resp_json))
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp_json)

}