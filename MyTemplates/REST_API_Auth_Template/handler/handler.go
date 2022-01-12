package handler

import (
	"fmt"
	"encoding/json"
	"strings"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"io"
	//"errors"
	//"strconv"
)

func typeCheck(data interface{}) interface{}{
	switch data.(type) {
    case string:
        d := data.(string)
        println("i is string:", d)

		dDate, err :=time.Parse("2006-01-02T15:04:05-07:00", d)
		if err != nil{
			fmt.Println(err)
			return d
		}
		utc, _ := time.LoadLocation("UTC")
		fmt.Println("i is Date %v",dDate)
		return dDate.In(utc)
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


func CreateFilter(baseData []map[string]interface{}) bson.D {
	filter := bson.D{}
	for _, data := range baseData {
		//keyのリストを取得
		keyList := []string{}
		for _, k := range data["key"].([]interface{}) {
			keyList = append(keyList, strings.ToLower(k.(string)))
		}
		filterKey := strings.Join(keyList, ".")
		fmt.Println(filterKey)
		//filter typeの取得
		filterType := data["type"].(float64)
		filterValue := typeCheck(data["value"])
		switch filterType {
		case 0: //完全一致
			filter = append(filter, bson.E{filterKey, filterValue})
		case 1: //リスト一致
			filter = append(filter, bson.E{filterKey, bson.D{{"$in", filterValue}} })
		case 2: //以上
			filter = append(filter, bson.E{filterKey, bson.D{{"$gte", filterValue}} })
		case 3: //以下
			filter = append(filter, bson.E{filterKey, bson.D{{"$lte", filterValue}} })
		case 4: //より上
			filter = append(filter, bson.E{filterKey, bson.D{{"$gt", filterValue}} })
		case 5: //未満
			filter = append(filter, bson.E{filterKey, bson.D{{"$lt", filterValue}} })
		}

	}
	fmt.Println("filter!!!")
	fmt.Println(filter)
	return filter
}

func BodyToFilter(body io.ReadCloser) (bson.D,error) {
	var getFilter []map[string]interface{}
	decoder := json.NewDecoder(body)
	//decoder.DisallowUnknownFields()
	err := decoder.Decode(&getFilter)
	if err != nil {
		return bson.D{}, err
	}
	filter := CreateFilter(getFilter)
	return filter, nil
}