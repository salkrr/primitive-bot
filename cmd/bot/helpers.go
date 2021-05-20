package main

import (
	"fmt"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/telegram"
)

func (app *application) serverError(chatID int64, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	err = app.errorLog.Output(2, trace)
	if err != nil {
		app.errorLog.Print(err)
	}

	_, err = app.bot.SendMessage(chatID, errorMessage)
	if err != nil {
		app.errorLog.Print(err)
	}
}

func (app *application) createStatusMessage(c primitive.Config, position int) string {
	return fmt.Sprintf(statusMessage, position, strings.ToLower(shapeNames[c.Shape]),
		c.Iterations, c.Repeat, c.Alpha, c.Extension, c.OutputSize)
}

func (app *application) getInputFromUser(
	chatID, menuMessageID int64,
	min, max int,
	in chan telegram.Message,
	out chan int,
) {
	err := app.bot.EditMessageText(chatID, menuMessageID,
		fmt.Sprintf(inputMessage, min, max))
	if err != nil {
		app.serverError(chatID, err)
		return
	}

	for {
		userMsg := <-in
		if err := app.bot.DeleteMessage(userMsg.Chat.ID, userMsg.MessageID); err != nil {
			app.serverError(chatID, err)
			return
		}

		userInput, err := strconv.Atoi(userMsg.Text)
		// correct input
		if err == nil && userInput >= min && userInput <= max {
			out <- userInput
			close(out)
			return
		}

		// incorrect input
		err = app.bot.EditMessageText(chatID, menuMessageID, fmt.Sprintf(inputMessage, min, max))
		if err != nil {
			if strings.Contains(err.Error(), "400") {
				// 400 error: message is not modified
				// and we don't care in this case
				continue
			}
			app.serverError(chatID, err)
			return
		}
	}
}

func (app *application) showMenu(
	chatID, messageID int64,
	chosenCallback,
	menuText string,
	template telegram.InlineKeyboardMarkup,
	params ...string,
) {
	newName := ""
	if len(params) > 0 {
		newName = params[0]
	}

	keyboard := newKeyboardFromTemplate(template, chosenCallback, newName)
	err := app.bot.EditMessageText(chatID, messageID, menuText, keyboard)
	if err != nil {
		if strings.Contains(err.Error(), "400") {
			// 400 error: message is not modified
			// and we don't care in this case
			return
		}
		app.serverError(chatID, err)
	}
}

// match reports whether path matches ^pattern$, and if it matches,
// assigns any capture groups to the *string or *int vars.
func match(path, pattern string, vars ...interface{}) bool {
	regex := mustCompileCached(fmt.Sprintf("^%s$", pattern))
	matches := regex.FindStringSubmatch(path)
	if len(matches) <= 0 {
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
