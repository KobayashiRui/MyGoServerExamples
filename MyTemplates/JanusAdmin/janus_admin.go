package main

import (
	"context"
	"log"
	"time"

	"github.com/rs/xid"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	websocketOptions := &websocket.DialOptions{}
	websocketOptions.Subprotocols = []string{"janus-protocol"}
	c, _, err := websocket.Dial(ctx, "ws://localhost:7188/admin", websocketOptions)
	if err != nil {
		log.Println(err.Error())
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	sendData := map[string]interface{}{}

	sendData["janus"] = "info"
	sendData["transaction"] = xid.New()
	err = wsjson.Write(ctx, c, &sendData)
	if err != nil {
		log.Println(err.Error())
	}

	readData := map[string]interface{}{}

	err = wsjson.Read(ctx, c, &readData)
	if err != nil {
		log.Println(err.Error())
	}

	log.Printf("%+v\n\n\n", readData)

	sendData = map[string]interface{}{}
	sendData["janus"] = "create"
	sendData["transaction"] = xid.New()

	err = wsjson.Write(ctx, c, &sendData)
	if err != nil {
		log.Println(err.Error())
	}

	readData = map[string]interface{}{}

	err = wsjson.Read(ctx, c, &readData)
	if err != nil {
		log.Println(err.Error())
	}

	log.Printf("%+v\n", readData)
	data := readData["data"].(map[string]interface{})
	sessionid := data["id"].(float64)
	log.Printf("sessionid : %v\n", sessionid)

	c.Close(websocket.StatusNormalClosure, "")
}
