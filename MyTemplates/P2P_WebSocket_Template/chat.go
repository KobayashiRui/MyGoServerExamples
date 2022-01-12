package main

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
	"fmt"
	"strings"
	"io/ioutil"
	"io"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
	"github.com/go-chi/chi/v5"
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
	rmsgs    chan []byte
	closeSlow func()
}
func (ss *socketServer)ConnectSocketHandlerHeader(w http.ResponseWriter, r *http.Request){
	websocketOptions := &websocket.AcceptOptions{}
	room :=  r.Header.Get("Sec-Websocket-Protocol")
	if room != "" {
		websocketOptions.Subprotocols = []string{room}
	}
	fmt.Printf("room:%v\n", room)
	//TODO TOKEN check
	c, err := websocket.Accept(w, r, websocketOptions)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	err = ss.connect(r.Context(), c, room)
	fmt.Println("END")
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

func (ss *socketServer)ConnectSocketHandler(w http.ResponseWriter, r *http.Request){
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")
	roomID := chi.URLParam(r, "roomID")
	fmt.Println(roomID)
	err = ss.connect(r.Context(), c, roomID)
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

func (ss *socketServer) SendHandler(w http.ResponseWriter, r *http.Request){
	roomID := chi.URLParam(r, "roomID")
	fmt.Println("SEND Handler")
	fmt.Println(roomID)
	body := http.MaxBytesReader(w, r.Body, 8192)
	msg, err := ioutil.ReadAll(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}

    ss.sendSocket(msg, roomID)
	w.Write([]byte("ok"))
	return
}

func (ss *socketServer) readMsg(ctx context.Context, c *websocket.Conn, s *socket) error {
	for{
		buf := new(strings.Builder)
		typ,r,err := c.Reader(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("type:%v\n", typ)
		len, err := io.Copy(buf, r)
		if err != nil {
			return fmt.Errorf("failed to io.Copy: %w", err)
		}
		if len > 0 {
			fmt.Printf("len : %v\n", len)
			fmt.Println(buf.String())
			s.rmsgs <- []byte(buf.String())
			fmt.Println("OK")
		}
	}
}

func (ss *socketServer) connect(ctx context.Context, c *websocket.Conn, room string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s:= &socket{
		msgs: make(chan []byte, ss.socketMessageBuffer),
		rmsgs: make(chan []byte, ss.socketMessageBuffer),
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}
	ss.addSocket(s, room)
	defer ss.deleteSocket(s, room)

	ctx_write, cancel_write := context.WithCancel(ctx)
	defer cancel_write()
	ticker := time.NewTicker(5*time.Second)
	defer ticker.Stop()

	ctx_read, cancel_read := context.WithCancel(ctx)
	defer cancel_read()
	go ss.readMsg(ctx_read, c, s)
	for {
    	select {
    	case msg := <-s.msgs:
    		err := writeTimeout(ctx_write, time.Second*5, c, msg)
    		if err != nil {
    		        return err
    		}
		case rmsg := <-s.rmsgs:
			fmt.Printf("read data %s\n", rmsg)
		case <-ticker.C:
			//fmt.Println("Ping");
			//cancel_read()
			//ctx_ping := c.CloseRead(ctx) // 1秒後にキャンセル
			//defer cancel_ping()
			pingErr := c.Ping(context.Background())
			if pingErr != nil {
				return pingErr
			}else{
				fmt.Println("get pong");
			}
			//go ss.readMsg(ctx_read, c, s)
		}

	}
}

func (ss *socketServer) sendSocket(msg []byte, room string) {
	ss.socketMu.Lock()
	defer ss.socketMu.Unlock()

	ss.socketLimiter.Wait(context.Background())

	for s := range ss.sockets[room] {
		select {
			case s.msgs <- msg:
			default:
					go s.closeSlow()
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
func (ss *socketServer) deleteSocket(s *socket, room string){
	fmt.Println("Delete: "+room)
	if len(ss.sockets[room]) != 0 {
           ss.socketMu.Lock()
           delete(ss.sockets[room], s)
           ss.socketMu.Unlock()
    }
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}