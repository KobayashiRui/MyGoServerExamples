package redisdb
import (
	"github.com/go-redis/redis/v8"
	"context"
	"time"
)

//セッション管理のredis関数

type SessionModel struct {
	Client *redis.Client
	Prefix string
}

//get user id  key = session id
// return user id
func (s *SessionModel) GetUserID(sessionId string) (string, error){
	ctx := context.Background()
	val, err := s.Client.Get(ctx, sessionId).Result()
	if err != nil {
		return "", err
	}
	return val, err
}

//set user id  key = session id
// return session id
func (s *SessionModel) SetUserID(sessionId string, userId string, exp time.Duration) error{
	ctx := context.Background()
	err := s.Client.Set(ctx, sessionId, userId, exp).Err()
	return err
}