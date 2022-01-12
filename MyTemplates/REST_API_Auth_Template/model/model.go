package model

import (
	"errors"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrNoRecord = errors.New("models: no matching record found") 

type User struct {
	ID  *primitive.ObjectID `json:"_id" bson:"_id"` 
	Email *string `json:",omitempty" bson:",omitempty"`
	Password *string `json:",omitempty" bson:",omitempty"`
	RegistDate *time.Time `json:",omitempty" bson:",omitempty"`
}

type TempUser struct {
	ID  *primitive.ObjectID `json:"_id" bson:"_id"` 
	Email *string `json:",omitempty" bson:",omitempty"`
	Password *string `json:",omitempty" bson:",omitempty"`
	Token *string `json:",omitempty" bson:",omitempty"`
	RegistDate *time.Time `json:",omitempty" bson:",omitempty"`
	//ExpireAt time.Time `json:"expire_at" bson:"expire_at,omitempty"` //有効期限を設定する場合
}
 func (u * User) SetFromTempUser(tu *TempUser) {
	if tu.ID != nil {
		u.ID = tu.ID
	}else{
		return
	}
	if tu.Email != nil {
		u.Email = tu.Email
	}else{
		return
	}
	if tu.Password != nil {
		u.Password = tu.Password
	}else{
		return
	}
	//fmt.Printf("OK")
 }
type Hoge struct {
	Name *string `json:",omitempty" bson:",omitempty"`
	A *float64 `json:",omitempty" bson:",omitempty"`
	B *int32 `json:",omitempty" bson:",omitempty"`
	C []int `json:",omitempty" bson:",omitempty"`
}

 type Article struct {
	ID  *primitive.ObjectID `json:"_id" bson:"_id"` 
	Name *string `json:",omitempty" bson:",omitempty"`
	A *float64 `json:",omitempty" bson:",omit"`
	B *int32 `json:",omitempty" bson:",omitempty"`
	C []int `json:",omitempty" bson:",omitempty"`
	RegistDate *time.Time `json:",omitempty" bson:",omitempty"`
	Hoges []Hoge `json:",omitempty" bson:",omitempty"`
 }