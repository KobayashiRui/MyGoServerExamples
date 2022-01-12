package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

// chatServer enables broadcasting to a set of subscribers.
type socketServer struct {
	socketMessageBuffer int
	socketLimiter       *rate.Limiter
	socketMu            sync.Mutex
	sockets             map[string]map[*socket]struct{}
	controllerSockets   map[string]map[*socket]struct{}
}

type socket struct {
	msgs      chan []byte
	rmsgs     chan []byte
	closeSlow func()
}

func NewSocketServer() socketServer {
	ss := socketServer{
		socketMessageBuffer: 16,
		sockets:             make(map[string]map[*socket]struct{}),
		socketLimiter:       rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}
	return ss
}

//websocketの接続
func (ss *socketServer) ConnectSocketHandler(w http.ResponseWriter, r *http.Request) {
	websocketOptions := &websocket.AcceptOptions{OriginPatterns: []string{"localhost:8080"}}
	//token := r.Header.Get("Sec-Websocket-Protocol")
	token := chi.URLParam(r, "tokenID")
	//if token != "" {
	//	websocketOptions.Subprotocols = []string{token}
	//}
	fmt.Printf("token:%v\n", token)
	//TODO TOKEN check
	c, err := websocket.Accept(w, r, websocketOptions)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	err = ss.connect(r.Context(), c, token)
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

func (ss *socketServer) SendIdHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "tokenID")
	//cameraID := chi.URLParam(r, "cameraID")
	//fmt.Println(cameraID)
	//token, err := ss.Camera.GetToken(cameraID)
	body := http.MaxBytesReader(w, r.Body, 8192)
	msg, err := ioutil.ReadAll(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}

	ss.sendSocket(msg, token)

	w.Write([]byte("OK"))
	return
}

func (ss *socketServer) readMsg(ctx context.Context, c *websocket.Conn, s *socket, token string) error {
	//ctx_read, cancel := context.WithCancel(ctx)
	//defer cancel()
	for {
		buf := new(strings.Builder)
		_, r, err := c.Reader(ctx)
		if err != nil {
			fmt.Println("Read ERROR 0")
			fmt.Println(err.Error())
			return err
		}

		//fmt.Println(t)

		len, err := io.Copy(buf, r)
		if err != nil {
			fmt.Println("Read ERROR 1")
			fmt.Println(err.Error())
			return err
			return fmt.Errorf("failed to io.Copy: %w", err)
		}
		if len > 0 {
			//fmt.Printf("len : %v\n", len)
			//fmt.Println(buf.String())
			for key := range ss.sockets[token] {
				key.rmsgs <- []byte(buf.String())
			}
			//s.rmsgs <- []byte("Hello")
		}
	}
}

//loop処理
func (ss *socketServer) connect(ctx context.Context, c *websocket.Conn, token string) error {
	//ctx, cancel := context.WithCancel(ctx)

	s := &socket{
		msgs:  make(chan []byte, ss.socketMessageBuffer),
		rmsgs: make(chan []byte, ss.socketMessageBuffer),
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}
	ss.addSocket(s, token)
	defer ss.deleteSocket(s, token)

	ctx_read, cancel := context.WithCancel(ctx)
	defer cancel()
	go ss.readMsg(ctx_read, c, s, token)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-s.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case rmsg := <-s.rmsgs:
			fmt.Printf("read data %v\n", rmsg)
			err := writeTimeout(ctx, time.Second*5, c, rmsg)
			if err != nil {
				return err
			}
		case <-ticker.C:
			//ctx_ping := c.CloseRead(ctx) // 1秒後にキャンセル
			pingErr := c.Ping(ctx)
			if pingErr != nil {
				fmt.Println("ping error")
				fmt.Println(pingErr.Error())
				return pingErr
			}
			//err := writeTimeout(ctx_read, time.Second*5, c, []byte("hello"))
			//if err != nil {
			//	return err
			//}
		}
	}
}
func (ss *socketServer) sendSocket(msg []byte, token string) {
	ss.socketMu.Lock()
	defer ss.socketMu.Unlock()

	ss.socketLimiter.Wait(context.Background())

	for s := range ss.sockets[token] {
		select {
		case s.msgs <- msg:
		default:
			go s.closeSlow()
		}
	}
}

func (ss *socketServer) addSocket(s *socket, token string) {
	ss.socketMu.Lock()
	if len(ss.sockets[token]) == 0 {
		ss.sockets[token] = make(map[*socket]struct{})
	}
	ss.sockets[token][s] = struct{}{}
	ss.socketMu.Unlock()
	//TODO Camera set online
}

func (ss *socketServer) deleteSocket(s *socket, token string) {
	fmt.Println("Delete: " + token)
	if len(ss.sockets[token]) != 0 {
		ss.socketMu.Lock()
		delete(ss.sockets[token], s)
		ss.socketMu.Unlock()
	}

}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}
