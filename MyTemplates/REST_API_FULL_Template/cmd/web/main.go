package main

import (
	"fmt"
	"context"
	"flag"
	"time"
	"net/http"
	"log"
	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
	"os"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"rest-api-temp-chi/dataservice/mongodb"
	"rest-api-temp-chi/handler"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	
	//MongoDB
	//uri := "mongodb://localhost:27017"
	uri := "mongodb://sampleAdmin:thisIsTest@localhost:27017"
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
        fmt.Println("connection error:", err)
    } else {
        fmt.Println("connection success:")
    }
	

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	
	h := &handler.Home{
		ErrorLog: errorLog, 
		InfoLog: infoLog,
	}
	ah := &handler.ArticleHandler{
		ErrorLog: errorLog, 
		InfoLog: infoLog,
		Article: &mongodb.ArticleModel{Client:client, DB:client.Database("test")},
	}

	router := chi.NewRouter()

	router.Get("/", h.Home)
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	//get article list
	//router.Get("/articles", )
	//get article data
	router.Get("/articles", ah.GetArticles)
	router.Get("/articles/{articleID}", ah.GetArticle)
	//create article
	router.Post("/articles", ah.CreateArticle)

	err = http.ListenAndServe(*addr, router)
	errorLog.Fatal(err)
}