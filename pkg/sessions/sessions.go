// Package sessions implements types for working with telegram sessions.
package sessions

import (
	"sync"
	"time"

	"github.com/lazy-void/primitive-bot/pkg/menu"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/tg"
)

//go:generate stringer -type=state
type state int

// Possible states of the user session.
const (
	InMenu state = iota
	InInputDialog
)

// Session represents one telegram session.
type Session struct {
	lastRequest   time.Time
	UserID        int64
	MenuMessageID int64
	State         state
	Input         chan tg.Message
	QuitInput     chan int
	ImgPath       string
	Menu          menu.Menu
	Config        primitive.Config
}

// NewSession initializes new instance of Session object.
func NewSession(userID, menuMessageID int64, imgPath string, workers int) Session {
	c := primitive.New(workers)

	return Session{
		lastRequest:   time.Now(),
		UserID:        userID,
		MenuMessageID: menuMessageID,
		State:         InMenu,
		Input:         make(chan tg.Message),
		QuitInput:     make(chan int),
		ImgPath:       imgPath,
		Menu:          menu.New(c),
		Config:        c,
	}
}

// ActiveSessions represents list of all active telegram sessions.
type ActiveSessions struct {
	sessions map[int64]Session
	timeout  time.Duration
	mu       sync.Mutex
}

// NewActiveSessions initializes new instance of ActiveSessions object.
// The argument 'timeout' specifies the maximum amount of time that a session
// can be inactive before it is terminated. The argument 'frequency' specifies
// how often the search for inactive sessions occurs.
func NewActiveSessions(timeout time.Duration, frequency time.Duration) *ActiveSessions {
	as := &ActiveSessions{
		sessions: make(map[int64]Session),
		timeout:  timeout,
	}
	go as.timeouter(frequency)

	return as
}

// Set adds new or replaces existing session.
func (as *ActiveSessions) Set(userID int64, s Session) {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.sessions[userID] = s
}

// Get returns session of user with specified ID. If the session
// doesn't exist, second parameter will be equal to false.
func (as *ActiveSessions) Get(userID int64) (Session, bool) {
	as.mu.Lock()
	defer as.mu.Unlock()

	s, ok := as.sessions[userID]
	if !ok {
		return Session{}, false
	}

	// Update info about time of last request
	s.lastRequest = time.Now()
	as.sessions[userID] = s

	return s, true
}

// timeouter terminates inactive sessions. The duration
// argument specifies interval between each search.
func (as *ActiveSessions) timeouter(d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		<-ticker.C

		as.mu.Lock()
		for _, s := range as.sessions {
			if time.Since(s.lastRequest) > as.timeout {
				if s.State == InInputDialog {
					s.QuitInput <- 1
				}

				delete(as.sessions, s.UserID)
			}
		}
		as.mu.Unlock()
	}
}
