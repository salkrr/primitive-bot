package sessions

import (
	"errors"
	"sync"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
)

var ErrNoSession = errors.New("session doesn't exist")

type Session struct {
	ChatID       int64
	ImgMessageID int64
	ImgPath      string
	Config       primitive.Config
}

type ActiveSessions struct {
	data map[int64]Session
	mu   sync.Mutex
}

func New() *ActiveSessions {
	return &ActiveSessions{
		data: make(map[int64]Session),
	}
}

func (as *ActiveSessions) Set(userID int64, s Session) {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.data[userID] = s
}

func (as *ActiveSessions) Get(userID int64) (Session, error) {
	as.mu.Lock()
	defer as.mu.Unlock()

	s, ok := as.data[userID]
	if !ok {
		return Session{}, ErrNoSession
	}

	return s, nil
}

func (as *ActiveSessions) Delete(userID int64) {
	as.mu.Lock()
	defer as.mu.Unlock()

	delete(as.data, userID)
}
