package mongodb
import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rest-api-temp-chi/model"
	"context"
	"time"
	"fmt"
)

type ArticleModel struct {
	Collect *mongo.Collection
}

func (m *ArticleModel) Get() ([]*model.Article, error) { 
	fmt.Println("GET Article")
	//collection := m.DB.Collection("Articles")
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//var doc bson.Raw
	//article := model.NewArticle()
	var results []*model.Article
	//article := new(model.Article)
	findOptions := options.Find()
	cur, err := m.Collect.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
	
	    // create a value into which the single document can be decoded
	    var article model.Article
	    err := cur.Decode(&article)
		if err != nil {
			return nil, err
		}
	    results = append(results, &article)
	}

	fmt.Println(results)
	return results,err
}

func (m *ArticleModel) GetOne(id string) (*model.Article, error) {
	fmt.Println("GET Article")
	//collection := m.DB.Collection("Articles")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := new(model.Article)
	object_id, err := primitive.ObjectIDFromHex(id)
	if(err != nil) {
		return nil, err
	}
	err = m.Collect.FindOne(ctx, bson.D{{"_id",object_id}}).Decode(result)
	fmt.Printf("%+v\n",*result)
	return result,err

}

func (m *ArticleModel) GetSearch(filter bson.D) ([]*model.Article, error) {
	fmt.Println("Search Article")
	//collection := m.DB.Collection("Articles")
	//err = m.Collect.FindOne(ctx, bson.D{{"_id",object_id}}).Decode(result)

	var results []*model.Article
	//article := new(model.Article)
	findOptions := options.Find()
	cur, err := m.Collect.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
	
	    // create a value into which the single document can be decoded
	    var article model.Article
	    err := cur.Decode(&article)
		if err != nil {
			return nil, err
		}
	    results = append(results, &article)
	}

	fmt.Println(results)
	return results,err

}

func (m *ArticleModel) Insert(newArticle model.Article) error {
	//collection := m.DB.Collection("Articles")

	new_id := primitive.NewObjectID()
	newArticle.ID = &new_id

	response, err := m.Collect.InsertOne(context.Background(), newArticle)
	fmt.Println(response)
	return err
}