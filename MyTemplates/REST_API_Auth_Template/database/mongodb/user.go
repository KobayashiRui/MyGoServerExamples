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

type UserModel struct {
	Collect *mongo.Collection
}

func (m *UserModel) Get() ([]*model.User, error) { 
	fmt.Println("GET User")
	//collection := m.DB.Collection("User")
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//var doc bson.Raw
	//article := model.NewArticle()
	var results []*model.User
	//article := new(model.Article)
	findOptions := options.Find()
	cur, err := m.Collect.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
	
	    // create a value into which the single document can be decoded
	    var user model.User
	    err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
	    results = append(results, &user)
	}

	fmt.Println(results)
	return results,err
}

func (m *UserModel) GetOne(id string) (*model.User, error) {
	fmt.Println("GET User")
	//collection := m.DB.Collection("User")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := new(model.User)
	object_id, err := primitive.ObjectIDFromHex(id)
	if(err != nil) {
		return nil, err
	}
	err = m.Collect.FindOne(ctx, bson.D{{"_id",object_id}}).Decode(result)
	fmt.Printf("%+v\n",*result)
	return result,err

}

func (m *UserModel) GetOneFromEmailPW(email string, password string) (*model.User, error) {
	fmt.Println("GET User")
	//collection := m.DB.Collection("User")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := new(model.User)
	err := m.Collect.FindOne(ctx, bson.D{{"email",email}, {"password",password}}).Decode(result)
	fmt.Printf("%+v\n",*result)
	return result,err

}
func (m *UserModel) GetOneFromEmail(email string) (*model.User, error) {
	fmt.Println("GET User from Email")
	//collection := m.DB.Collection("User")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := new(model.User)
	filter := bson.D{{"email",email}}
	err := m.Collect.FindOne(ctx, filter).Decode(result)
	fmt.Println(err)
	fmt.Printf("%+v\n",*result)
	return result,err

}

func (m *UserModel) Insert(newUser model.User) error {
	//collection := m.DB.Collection("User")

	new_id := primitive.NewObjectID()
	newUser.ID = &new_id
	//user := model.User{
	//	ID: primitive.NewObjectID(),
	//	Email: email,
	//	Password: password,
	//	Authority: 0,
	//}
	response, err := m.Collect.InsertOne(context.Background(), newUser)
	fmt.Println(response)
	return err
}