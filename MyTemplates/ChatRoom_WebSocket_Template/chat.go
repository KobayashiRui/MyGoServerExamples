package main

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
	"fmt"
	"strings"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

// chatServer enables broadcasting to a set of subscribers.
type socketServer struct {
	socketMessageBuffer int
	socketLimiter *rate.Limiter
	socketMu sync.Mutex
	sockets   map[string]map[*socket]struct{}
}

type socket struct {
	msgs      chan []byte
	closeSlow func()
}

func (ss *socketServer)connectSocketHandler(w http.ResponseWriter, r *http.Request){
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	room := strings.Split(r.URL.Path, "/")
	fmt.Println(room)
	err = ss.connect(r.Context(), c, room[2])
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

}


func (ss *socketServer) connect(ctx context.Context, c *websocket.Conn, room string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s:= &socket{
		msgs: make(chan []byte, ss.socketMessageBuffer),
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}
	ss.addSocket(s, room)
	for {
		fmt.Println("read")
		typ, r, err := c.Read(ctx)
		fmt.Println(typ)
		fmt.Println(string(r))		
		if err != nil {
			fmt.Println("err: ", err)
			return err
		}
		receivedMessage := string(r)
		err = c.Write(ctx, websocket.MessageText, []byte("hogehoge-"+receivedMessage))
		if err != nil {
			fmt.Println("err: ", err)
			return err
		}
	}
}


func (ss *socketServer) addSocket(s *socket, room string) {
	ss.socketMu.Lock()
	if len(ss.sockets[room]) == 0 {
		ss.sockets[room] = make(map[*socket]struct{})
	  }
	ss.sockets[room][s] = struct{}{}
	ss.socketMu.Unlock()
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}