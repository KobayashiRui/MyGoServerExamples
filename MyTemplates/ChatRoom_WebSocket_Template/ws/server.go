package ws

type socketServer struct {
	socketMessageBuffer int
	socketLimiter *rate.Limiter
	socketMu sync.Mutex
	sockets   map[string]map[*socket]struct{}
}

