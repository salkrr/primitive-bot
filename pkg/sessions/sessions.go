package sessions

import (
	"sync"

	"github.com/lazy-void/primitive-bot/pkg/menu"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/telegram"
)

type Session struct {
	ChatID        int64
	MenuMessageID int64
	InChan        chan<- telegram.Message
	ImgPath       string
	Menu          menu.Menu
	Config        primitive.Config
}

func NewSession(chatID, menuMessageID int64, imgPath string) Session {
	c := primitive.NewConfig()

	return Session{
		ChatID:        chatID,
		MenuMessageID: menuMessageID,
		ImgPath:       imgPath,
		Menu:          menu.New(c),
		Config:        c,
	}
}

type ActiveSessions struct {
	data map[int64]Session
	mu   sync.Mutex
}

func NewActiveSessions() *ActiveSessions {
	return &ActiveSessions{
		data: make(map[int64]Session),
	}
}

func (as *ActiveSessions) Set(userID int64, s Session) {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.data[userID] = s
}

func (as *ActiveSessions) Get(userID int64) (Session, bool) {
	as.mu.Lock()
	defer as.mu.Unlock()

	s, ok := as.data[userID]
	if !ok {
		return Session{}, false
	}

	return s, true
}

func (as *ActiveSessions) Delete(userID int64) {
	as.mu.Lock()
	defer as.mu.Unlock()

	delete(as.data, userID)
}
