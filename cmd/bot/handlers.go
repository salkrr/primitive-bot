package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/queue"
	"github.com/lazy-void/primitive-bot/pkg/sessions"
	"github.com/lazy-void/primitive-bot/pkg/telegram"
)

func (app *application) handleMessage(m telegram.Message) {
	if m.Photo == nil {
		// Handle user input if they are inside input form
		s, ok := app.sessions.Get(m.From.ID)
		if ok && s.InChan != nil {
			s.InChan <- m
			return
		}

		// Send help message
		err := app.bot.SendMessage(m.Chat.ID, helpMessage)
		if err != nil {
			app.serverError(m.Chat.ID, err)
		}
		return
	}

	// If we already have session - delete it's menu
	s, ok := app.sessions.Get(m.From.ID)
	if ok {
		err := app.bot.DeleteMessage(s.ChatID, s.MenuMessageID)
		if err != nil {
			app.serverError(m.Chat.ID, err)
			return
		}
	}

	// Choose smallest image with dimensions >= 256
	var file telegram.PhotoSize
	for _, photo := range m.Photo {
		file = photo
		if photo.Width >= 256 && photo.Height >= 256 {
			break
		}
	}
	if file.FileID == "" {
		app.serverError(m.Chat.ID, fmt.Errorf("no image files in %v", m.Photo))
		return
	}

	path := fmt.Sprintf("%s/%s.jpg", app.inDir, file.FileUniqueID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		img, err := app.bot.DownloadFile(file.FileID)
		if err != nil {
			app.serverError(m.Chat.ID, fmt.Errorf("couldn't download image: %s", err))
			return
		}
		os.WriteFile(path, img, 0644)
	}

	// Create session
	msg, err := app.bot.SendMessageWithInlineKeyboard(m.Chat.ID, rootMenuText, rootKeyboard)
	if err != nil {
		app.serverError(m.Chat.ID, err)
		return
	}
	app.sessions.Set(m.From.ID, sessions.Session{
		ChatID:        m.Chat.ID,
		MenuMessageID: msg.MessageID,
		InChan:        nil,
		ImgPath:       path,
		Config:        primitive.NewConfig(),
	})
}

func (app *application) handleCallbackQuery(q telegram.CallbackQuery) {
	defer func() {
		if err := app.bot.AnswerCallbackQuery(q.ID, ""); err != nil {
			app.errorLog.Printf("Answer callback query error: %s", err)
		}
	}()

	s, ok := app.sessions.Get(q.From.ID)
	if !ok || q.Message.MessageID != s.MenuMessageID {
		err := app.bot.DeleteMessage(q.Message.Chat.ID, q.Message.MessageID)
		if err != nil {
			app.errorLog.Printf("Deleting message error: %s", err)
		}
		return
	}

	var num int
	var slug string
	switch {
	case match(q.Data, "/"):
		app.bot.EditMessageTextWithKeyboard(q.Message.Chat.ID, q.Message.MessageID, rootMenuText, rootKeyboard)
	case match(q.Data, "/start"):
		// app.sessions.Delete(q.From.ID)
		// err := app.bot.DeleteMessage(q.Message.Chat.ID, q.Message.MessageID)
		// if err != nil {
		// 	app.serverError(q.Message.Chat.ID, err)
		// }

		pos := app.queue.Enqueue(queue.Operation{
			ChatID:  q.Message.Chat.ID,
			ImgPath: s.ImgPath,
			Config:  s.Config,
		})
		report := fmt.Sprintf(enqueuedMessage, pos, strings.ToLower(shapeNames[s.Config.Shape]),
			s.Config.Iterations, s.Config.Repeat, s.Config.Alpha, s.Config.Extension, s.Config.OutputSize)
		err := app.bot.SendMessage(q.Message.Chat.ID, report)
		if err != nil {
			app.serverError(q.Message.Chat.ID, err)
		}
	case match(q.Data, "/settings/shape"):
		optionCallback := fmt.Sprintf("/settings/shape/%d", s.Config.Shape)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, optionCallback, shapesMenuText, shapesKeyboard)
	case match(q.Data, "/settings/shape/([0-8])", &num):
		s.Config.Shape = primitive.Shape(num)
		app.sessions.Set(q.From.ID, s)

		// update menu
		optionCallback := fmt.Sprintf("/settings/shape/%d", s.Config.Shape)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, optionCallback, shapesMenuText, shapesKeyboard)
	case match(q.Data, "/settings/iter"):
		optionCallback := fmt.Sprintf("/settings/iter/%d", s.Config.Iterations)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, optionCallback, iterMenuText, iterKeyboard)
	case match(q.Data, "/settings/iter/([0-9]+)", &num):
		if num > 10000 {
			return
		}
		s.Config.Iterations = num
		app.sessions.Set(q.From.ID, s)

		// update menu
		optionCallback := fmt.Sprintf("/settings/iter/%d", s.Config.Iterations)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, optionCallback, iterMenuText, iterKeyboard)
	case match(q.Data, "/settings/iter/input"):
		app.handleIterInput(q, s)
	case match(q.Data, "/settings/rep"):
		optionCallback := fmt.Sprintf("/settings/rep/%d", s.Config.Repeat)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, optionCallback, repMenuText, repKeyboard)
	case match(q.Data, "/settings/rep/([1-6])", &num):
		s.Config.Repeat = num
		app.sessions.Set(q.From.ID, s)

		// update menu
		optionCallback := fmt.Sprintf("/settings/rep/%d", s.Config.Repeat)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, optionCallback, repMenuText, repKeyboard)
	case match(q.Data, "/settings/alpha"):
		optionCallback := fmt.Sprintf("/settings/alpha/%d", s.Config.Alpha)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, optionCallback, alphaMenuText, alphaKeyboard)
	case match(q.Data, "/settings/alpha/([0-9]+)", &num):
		if num < 0 || num > 255 {
			return
		}
		s.Config.Alpha = num
		app.sessions.Set(q.From.ID, s)

		// update menu
		optionCallback := fmt.Sprintf("/settings/alpha/%d", s.Config.Alpha)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, optionCallback, alphaMenuText, alphaKeyboard)
	case match(q.Data, "/settings/alpha/input"):
		app.handleAlphaInput(q, s)
	case match(q.Data, "/settings/ext"):
		optionCallback := fmt.Sprintf("/settings/ext/%s", s.Config.Extension)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, optionCallback, extMenuText, extKeyboard)
	case match(q.Data, "/settings/ext/(jpg|png|svg|gif)", &slug):
		s.Config.Extension = primitive.Extension(slug)
		app.sessions.Set(q.From.ID, s)

		// update menu
		optionCallback := fmt.Sprintf("/settings/ext/%s", s.Config.Extension)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, optionCallback, extMenuText, extKeyboard)
	case match(q.Data, "/settings/size"):
		optionCallback := fmt.Sprintf("/settings/size/%d", s.Config.OutputSize)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, optionCallback, sizeMenuText, sizeKeyboard)
	case match(q.Data, "/settings/size/([0-9]+)", &num):
		if num < 256 || num > 1920 {
			return
		}
		s.Config.OutputSize = num
		app.sessions.Set(q.From.ID, s)

		// update menu
		optionCallback := fmt.Sprintf("/settings/size/%d", s.Config.OutputSize)
		app.showMenu(q.Message.Chat.ID, q.Message.MessageID, optionCallback, sizeMenuText, sizeKeyboard)
	case match(q.Data, "/settings/size/input"):
		app.handleSizeInput(q, s)
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
	err := app.bot.EditMessageTextWithKeyboard(chatID, messageID, menuText, keyboard)
	if err != nil {
		if strings.Contains(err.Error(), "400") {
			// 400 error: message is not modified
			// and we don't care in this case
			return
		}
		app.serverError(chatID, err)
	}
}

func (app *application) handleIterInput(
	q telegram.CallbackQuery,
	s sessions.Session,
) {
	in := make(chan telegram.Message)
	out := make(chan int)

	// Get user input
	s.InChan = in
	app.sessions.Set(q.From.ID, s)

	go app.getInputFromUser(q.Message.Chat.ID, q.Message.MessageID, 1, 5000, in, out)

	s.Config.Iterations = <-out
	s.InChan = nil
	app.sessions.Set(q.From.ID, s)
	close(in)

	buttonText := fmt.Sprintf("Другое (%d)", s.Config.Iterations)
	app.showMenu(q.Message.Chat.ID, q.Message.MessageID, "/settings/iter/input",
		iterMenuText, iterKeyboard, buttonText)
}

func (app *application) handleAlphaInput(
	q telegram.CallbackQuery,
	s sessions.Session,
) {
	in := make(chan telegram.Message)
	out := make(chan int)

	// Get user input
	s.InChan = in
	app.sessions.Set(q.From.ID, s)

	go app.getInputFromUser(q.Message.Chat.ID, q.Message.MessageID, 1, 255, in, out)

	s.Config.Alpha = <-out
	s.InChan = nil
	app.sessions.Set(q.From.ID, s)
	close(in)

	buttonText := fmt.Sprintf("Другое (%d)", s.Config.Alpha)
	app.showMenu(q.Message.Chat.ID, q.Message.MessageID, "/settings/alpha/input",
		alphaMenuText, alphaKeyboard, buttonText)
}

func (app *application) handleSizeInput(
	q telegram.CallbackQuery,
	s sessions.Session,
) {
	in := make(chan telegram.Message)
	out := make(chan int)

	// Get user input
	s.InChan = in
	app.sessions.Set(q.From.ID, s)

	go app.getInputFromUser(q.Message.Chat.ID, q.Message.MessageID, 256, 3840, in, out)

	s.Config.OutputSize = <-out
	s.InChan = nil
	app.sessions.Set(q.From.ID, s)
	close(in)

	buttonText := fmt.Sprintf("Другое (%d)", s.Config.OutputSize)
	app.showMenu(q.Message.Chat.ID, q.Message.MessageID, "/settings/size/input",
		sizeMenuText, sizeKeyboard, buttonText)
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

func (app *application) worker() {
	for {
		// get next operation
		op, ok := app.queue.Peek()
		if !ok {
			continue
		}

		start := time.Now()
		app.infoLog.Printf("Creating from '%s' for chat '%d': count=%d, mode=%d, alpha=%d, repeat=%d, resolution=%d, extension=%s",
			op.ImgPath, op.ChatID, op.Config.Iterations, op.Config.Shape, op.Config.Alpha, op.Config.Repeat, op.Config.OutputSize, op.Config.Extension)

		// create primitive
		outputPath := fmt.Sprintf("%s/%d_%d.%s", app.outDir, op.ChatID, start.Unix(), op.Config.Extension)
		err := op.Config.Create(op.ImgPath, outputPath)
		if err != nil {
			app.serverError(op.ChatID, err)
			return
		}
		app.infoLog.Printf("Finished creating '%s' for chat '%d'; Output: '%s'; Time: %.1f seconds",
			filepath.Base(op.ImgPath), op.ChatID, filepath.Base(outputPath), time.Since(start).Seconds())

		// send output to the user
		err = app.bot.SendDocument(op.ChatID, outputPath)
		if err != nil {
			app.serverError(op.ChatID, err)
			return
		}
		app.infoLog.Printf("Sent result '%s' to the chat '%d'", filepath.Base(outputPath), op.ChatID)

		// remove operation from the queue
		app.queue.Dequeue()
	}
}
