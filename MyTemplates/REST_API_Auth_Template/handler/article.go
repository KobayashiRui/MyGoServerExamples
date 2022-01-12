package handler
import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"rest-api-temp-chi/database/mongodb"
	"rest-api-temp-chi/model"
	"time"
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
		w.WriteHeader(http.StatusBadRequest)
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

	article_json, err := json.Marshal(article)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ah.InfoLog.Println(string(article_json))
	w.Header().Set("Content-Type", "application/json")
	w.Write(article_json)
}

func (ah *ArticleHandler) SearchArticles(w http.ResponseWriter, r *http.Request){

	//var getFilter []map[string]interface{}
	//decoder := json.NewDecoder(r.Body)
	//decoder.DisallowUnknownFields()
	//err := decoder.Decode(&getFilter)
	//filter := CreateFilter(getFilter)
	filter, err := BodyToFilter(r.Body)
	println("Search")
	//println(getFilter)
	articles, err_db := ah.Article.GetSearch(filter)
	//fmt.Println(*article)
	if err_db != nil {
		ah.ErrorLog.Println(err_db.Error())
		//serverError(w, err_db)
		return
	}

	article_json, err := json.Marshal(articles)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ah.InfoLog.Println(string(article_json))
	w.Header().Set("Content-Type", "application/json")
	w.Write(article_json)
}


func (ah *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {

	var new_article model.Article
	decoder := json.NewDecoder(r.Body)
	//decoder.DisallowUnknownFields()
	err := decoder.Decode(&new_article)
	fmt.Printf("%+v\n", new_article) // with Variable name
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	nowTime := time.Now()
	new_article.RegistDate = &nowTime
	err = ah.Article.Insert(new_article)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte("ok"))
	return
}