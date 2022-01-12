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
	"github.com/go-redis/redis/v8"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"rest-api-temp-chi/database/mongodb"
	"rest-api-temp-chi/database/redisdb"
	"rest-api-temp-chi/auth"
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

	//redisdb
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	_userDB := &mongodb.UserModel{Collect: client.Database("test").Collection("User")}
	_tempUserDB := &mongodb.TempUserModel{Collect: client.Database("test").Collection("TempUser")}
	_articleDB := &mongodb.ArticleModel{Collect: client.Database("test").Collection("Article")}
	
	ah := &auth.AuthHandler{
		ErrorLog: errorLog, 
		InfoLog: infoLog,
		TempUser: _tempUserDB,
		User: _userDB,
		Session:  &redisdb.SessionModel{Client: rdb, Prefix: ""},
	}

	h := &handler.ArticleHandler{
		ErrorLog: errorLog, 
		InfoLog: infoLog,
		Article: _articleDB,
	}


	router := chi.NewRouter()
	router.Post("/signup", ah.SignUp)
	router.Get("/confirm/{userID}", ah.ConfirmSignUp)
	router.Post("/signin", ah.SignIn)
	router.Get("/tempusers", ah.GetTempUsers)
	router.Get("/users", ah.GetUsers)

	router.Post("/test/search", handler.TestSearch)

	router.Route("/welcome", func(r chi.Router) {
		r.Use(ah.Auth)
		r.Get("/", handler.Welcome)
	})

	router.Get("/article", h.GetArticles)
	router.Post("/article", h.CreateArticle)
	router.Post("/article/search", h.SearchArticles)


	err = http.ListenAndServe(*addr, router)
	errorLog.Fatal(err)
}