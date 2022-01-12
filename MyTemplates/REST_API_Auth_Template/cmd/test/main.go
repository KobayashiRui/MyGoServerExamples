package main

import (
	"time"
	"fmt"
	"encoding/json"
	"reflect"
	"strings"
	"go.mongodb.org/mongo-driver/bson"
)


type Test struct{
	Name *string `json:",omitempty" bson:",omitempty"`
	Age *int `json:",omitempty" bson:",omitempty"`
	RegistDate *time.Time `json:",omitempty" bson:",omitempty"`
}

func CreateFilter(baseData []map[string]interface{}){
	filter := bson.D{}
	for _, data := range baseData {
		//keyのリストを取得
		keyList := []string{}
		for _, k := range data["key"].([]interface{}) {
			keyList = append(keyList, k.(string))
		}
		filterKey := strings.Join(keyList, ".")
		fmt.Printf(filterKey)
		//filter typeの取得
		filterType := data["type"].(float64)
		switch filterType {
		case 0: //完全一致
			filter = append(filter, bson.E{filterKey, data["value"]})
		case 1: //リスト一致
		case 2: //以上
		case 3: //以下
		case 4: //より上
		case 5: //未満
		}

	}
	fmt.Println("filter!!!")
	fmt.Println(filter)
}


func TypeCheck(data interface{}) interface{}{
	switch data.(type) {
    case string:
        d := data.(string)
        println("i is string:", d)

		dDate, err :=time.Parse("2006-01-02T15:04:05-07:00", d)
		if err != nil{
			return d
		}
		return dDate
    case bool:
        d := data.(bool)
        println("i is boolean:", d)
		return d
    case int:
        d := data.(int)
        println("i is integer:", d) // i is integer: 100
		return d
	case float64:
		d := data.(float64)
		println("i is float64:", d) // i is float64
		return d
	case []interface{}:
		println("array")
		d := data.([]interface{})
		//arrayType := typeCheck(d[0])
		return d
	default:
		println("ERROR")
		return 0
	}
}

/*
func typeCheck(data interface{}){
	r := reflect.ValueOf(data)
	if r.IsValid() {
        switch r.Kind() {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			println("int")
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			println("uint")
        case reflect.Float32, reflect.Float64:
			println("float")
		case reflect.Slice:
			println("arrray")
			typeCheck()
		case reflect.String:
			println("string")
        default:
			println(r.Kind())
        }
    }
}*/

func parseFilter(data map[string]interface{}, baseStruct interface{}) {
	fmt.Printf("data:%+v\n", data);
	fmt.Printf("%+v\n", data["key"]);
	//reflect.ValueOf(hoge).Elem().FieldByName().Type()
	fmt.Println("types")
	//TypeCheck(data["key"])
	//TypeCheck(data["value"])
	keyList := []string{}
	for i, k := range data["key"].([]interface{}) {
    		fmt.Println(i, k)
		keyList = append(keyList, k.(string))
		
	}
	fmt.Printf("%+v\n",keyList)
	
	typeData := uint8(data["type"].(float64)) //0~5
	fmt.Printf("%v\n",typeData)
	//TODO エラー処理
	typeValue := reflect.ValueOf(&baseStruct).Elem().Elem().FieldByName(keyList[0]).Type()
	fmt.Printf("Value type:%+v\n",typeValue)
}

func main() {
myJSON := `[
	{"key": ["Name", "Fuga"],
	 "value": "Unko",
	 "type":0
	},
	{"key": ["Age"],
	 "value": [10,20,50],
	 "type":2
	},
	{"key": ["RegistDate"],
	 "value": "2021-11-02T09:00:00+09:00",
	 "type":0
	}]`
//tData, err :=time.Parse("2006-01-02T15:04:05-07:00", "2021-11-02T09:00:00+09:00")
tData, err :=time.Parse("2006-01-02T15:04:05", "2021-11-02T09:00:00+09:00")
if err != nil{
fmt.Println(err)
}
fmt.Println(tData)
bsonTest := bson.D{}
bsonTest = append(bsonTest, bson.E{strings.ToLower("BaseData"), "Fuga"})
bsonTest = append(bsonTest, bson.E{strings.ToLower("Date"), time.Now()})
bsonTest = append(bsonTest, bson.E{strings.ToLower("Date2"), tData})

bsonD, err := bson.Marshal(bsonTest)
if err != nil{
	fmt.Println(err)
}

fmt.Println("bsondata: %s\n",string(bsonD))

fmt.Println("ok")
var result []map[string]interface{}
json.Unmarshal([]byte(myJSON), &result)

fmt.Printf("reuslt : %+v\n", result)
for i, s := range result {
    fmt.Printf("index: %d, %+v\n", i, s)
    parseFilter(s, Test{})
}

CreateFilter(result)
}