package main

import (
	"fmt"
	"runtime/debug"
)

func (app *application) serverError(chatID int64, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)
	app.bot.SendMessage(chatID, messageError)
}
