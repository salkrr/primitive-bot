package main

import (
	"errors"
	"fmt"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/lazy-void/primitive-bot/pkg/menu"
	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/sessions"
)

var errSessionTerminated = errors.New("session terminated")

func (app *application) serverError(chatID int64, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	err = app.errorLog.Output(2, trace)
	if err != nil {
		app.errorLog.Print(err)
	}

	_, err = app.bot.SendMessage(chatID,
		app.printer.Sprintf("Something gone wrong! Please, try again in a few minutes."))
	if err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) createStatusMessage(c primitive.Config, position int) string {
	return app.printer.Sprintf(
		"%d place in the queue.\n\nShapes: %s\nSteps: %d\nRepetitions: %d\nAlpha-channel: %d\nExtension: %s\nSize: %#v",
		position, strings.ToLower(menu.ShapeNames[c.Shape]), c.Iterations, c.Repeat, c.Alpha, c.Extension, c.OutputSize,
	)
}

func (app *application) getInputFromUser(
	s sessions.Session,
	min, max int,
) (int, error) {
	s.State = sessions.InInputDialog
	app.sessions.Set(s.UserID, s, false)

	err := app.bot.EditMessageText(s.UserID, s.MenuMessageID,
		app.printer.Sprintf("Enter number between %#v and %#v:", min, max))
	if err != nil {
		return 0, err
	}

	for {
		select {
		case msg := <-s.Input:
			// Delete message with user input
			err := app.bot.DeleteMessage(msg.Chat.ID, msg.MessageID)
			if err != nil {
				return 0, err
			}

			userInput, err := strconv.Atoi(msg.Text)
			// correct input
			if err == nil && userInput >= min && userInput <= max {
				s.State = sessions.InMenu
				return userInput, nil
			}

			// incorrect input
			err = app.bot.EditMessageText(
				s.UserID, s.MenuMessageID,
				app.printer.Sprintf("Incorrect value!\nEnter number between %#v and %#v:", min, max),
			)
			if err != nil {
				if strings.Contains(err.Error(), "400") {
					// 400 error: message is not modified
					// and we don't care in this case
					break
				}
				return 0, err
			}
		case <-s.QuitInput:
			return 0, errSessionTerminated
		}
	}
}

func (app *application) showMenuView(
	chatID, messageID int64,
	view menu.View,
) {
	err := app.bot.EditMessageText(chatID, messageID, view.Text, view.Keyboard)
	if err != nil {
		if strings.Contains(err.Error(), "400") {
			// 400 error: message is not modified
			// and we don't care in this case
			return
		}
		app.serverError(chatID, err)
	}
}

func (app *application) sendMessage(chatID int64, message string) {
	_, err := app.bot.SendMessage(chatID, message)
	if err != nil {
		app.serverError(chatID, err)
	}
}

// match reports whether path matches ^pattern$, and if it matches,
// assigns any capture groups to the *string or *int vars.
func match(path, pattern string, vars ...interface{}) bool {
	regex := mustCompileCached(fmt.Sprintf("^%s$", pattern))
	matches := regex.FindStringSubmatch(path)
	if len(matches) == 0 {
		return false
	}

	for i, match := range matches[1:] {
		switch p := vars[i].(type) {
		case *string:
			*p = match
		case *int:
			n, err := strconv.Atoi(match)
			if err != nil {
				return false
			}
			*p = n
		default:
			return false
		}
	}
	return true
}

var (
	regexen = make(map[string]*regexp.Regexp)
	relock  = sync.Mutex{}
)

func mustCompileCached(pattern string) *regexp.Regexp {
	relock.Lock()
	defer relock.Unlock()

	regex := regexen[pattern]
	if regex == nil {
		regex = regexp.MustCompile(pattern)
		regexen[pattern] = regex
	}
	return regex
}
