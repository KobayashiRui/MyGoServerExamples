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

type TempUserModel struct {
	Collect *mongo.Collection
}

func (m *TempUserModel) Get() ([]*model.TempUser, error) { 
	fmt.Println("GET User")
	//collection := m.DB.Collection("TempUser")
	var results []*model.TempUser
	findOptions := options.Find()
	cur, err := m.Collect.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
	
	    // create a value into which the single document can be decoded
	    var user model.TempUser
	    err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
	    results = append(results, &user)
	}

	fmt.Println(results)
	return results,err
}

//idからユーザー情報を取得
func (m *TempUserModel) GetOne(id string) (*model.TempUser, error) {
	fmt.Println("GET User")
	//collection := m.DB.Collection("TempUser")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := new(model.TempUser)
	object_id, err := primitive.ObjectIDFromHex(id)
	if(err != nil) {
		return nil, err
	}
	err = m.Collect.FindOne(ctx, bson.D{{"_id",object_id}}).Decode(result)
	fmt.Printf("%+v\n",*result)
	return result,err
}

func (m *TempUserModel) GetOneFromToken(token string) (*model.TempUser, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := new(model.TempUser)
	err := m.Collect.FindOne(ctx, bson.D{{"Token",token}}).Decode(result)
	return result,err
}

func (m *TempUserModel) Insert(newTempUser model.TempUser) (*model.TempUser, error) {
	//collection := m.DB.Collection("TempUser")
	new_id := primitive.NewObjectID()
	newTempUser.ID = &new_id
	nowTime := time.Now()
	newTempUser.RegistDate = &nowTime

	bsonData, err := bson.Marshal(newTempUser)
	fmt.Println("newTempUser: %+v\n", bsonData)
	res, err := m.Collect.InsertOne(context.Background(), newTempUser)
	fmt.Println(res)
	return &newTempUser, err
}


func (m *TempUserModel) DeleteOne(id string) error {

	//collection := cc.DB.Collection("ProcedureSocialInsRelated")
	object_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := m.Collect.DeleteOne(context.Background(), bson.D{{"_id",object_id}})
	if err != nil {
	    fmt.Println(err)
	}
	fmt.Printf("DeleteOne removed %v document(s)\n", res.DeletedCount)
	return err
}

//mongoでTTL(存続時間)を設定する場合に使用する.
//func (m *TempUserModel) SetExpire() {
//	//mongo expire set 
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	expireIndex := mongo.IndexModel{ 
//		Keys: bson.D{{Key: "expire_at", Value: 1}}, 
//		Options: options.Index().SetExpireAfterSeconds(0), 
//	}
//	m.Collect.Indexes().CreateOne(ctx, expireIndex)
//}