package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/sessions"
	"github.com/lazy-void/primitive-bot/pkg/telegram"
)

func (app *application) handleMessage(m telegram.Message, out chan sessions.Session) {
	if m.Photo == nil {
		// Handle user input if they are inside input form
		s, ok := app.sessions.Get(m.From.ID)
		if ok && s.InputChannel != nil {
			s.InputChannel <- m
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
		if photo.Width >= 256 && photo.Height >= 256 {
			file = photo
			break
		}
	}
	if file.FileID == "" {
		app.serverError(m.Chat.ID, fmt.Errorf("no image with dimensions >= 256 in %v", m.Photo))
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
	msg, err := app.bot.SendMessageWithInlineKeyboard(m.Chat.ID, rootMenu, rootKeyboard)
	if err != nil {
		app.serverError(m.Chat.ID, err)
		return
	}
	app.sessions.Set(m.From.ID, sessions.Session{
		ChatID:        m.Chat.ID,
		MenuMessageID: msg.MessageID,
		InputChannel:  nil,
		ImgPath:       path,
		Config:        primitive.NewConfig(),
	})
}

func (app *application) handleCallbackQuery(q telegram.CallbackQuery, out chan sessions.Session) {
	defer app.bot.AnswerCallbackQuery(q.ID, "")

	s, ok := app.sessions.Get(q.From.ID)
	if !ok {
		// Remove keyboard of an invalid session
		app.bot.DeleteMessage(q.Message.Chat.ID, q.Message.MessageID)
		return
	}

	var num int
	var slug string
	switch {
	case match(q.Data, "/"):
		app.bot.EditMessageTextWithKeyboard(q.Message.Chat.ID, q.Message.MessageID, rootMenu, rootKeyboard)
	case match(q.Data, "/start"):
		app.sessions.Delete(q.From.ID)
		app.bot.DeleteMessage(q.Message.Chat.ID, q.Message.MessageID)
		app.bot.SendMessage(q.Message.Chat.ID, fmt.Sprintf(enqueuedMessage,
			strings.ToLower(shapeNames[s.Config.Shape]), s.Config.Iterations,
			s.Config.Repeat, s.Config.Alpha, s.Config.Extension, s.Config.OutputSize))

		out <- s
	case match(q.Data, "/settings/shape"):
		app.bot.EditMessageTextWithKeyboard(q.Message.Chat.ID, q.Message.MessageID, shapesMenu, shapesKeyboard)
	case match(q.Data, "/settings/shape/([0-8])", &num):
		// TODO: add symbol to the chosen option
		s.Config.Shape = primitive.Shape(num)
		app.sessions.Set(q.From.ID, s)
		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Будут использоваться фигуры: %s.", strings.ToLower(shapeNames[primitive.Shape(num)])))
	case match(q.Data, "/settings/iter"):
		app.bot.EditMessageTextWithKeyboard(q.Message.Chat.ID, q.Message.MessageID, iterMenu, iterKeyboard)
	case match(q.Data, "/settings/iter/([0-9]+)", &num):
		// TODO: add symbol to the chosen option
		if num > 10000 {
			return
		}

		s.Config.Iterations = num
		app.sessions.Set(q.From.ID, s)
		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Поменял количество итераций на %d.", num))
	case match(q.Data, "/settings/iter/diff"):
		in := make(chan telegram.Message)
		out := make(chan int)

		// Get user input
		s.InputChannel = in
		app.sessions.Set(q.From.ID, s)
		go app.getInputFromUser(
			q.Message.Chat.ID, q.Message.MessageID, 1, 5000,
			iterMenu, iterKeyboard, in, out,
		)

		s.Config.Iterations = <-out
		s.InputChannel = nil
		app.sessions.Set(q.From.ID, s)
		close(in)

		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Поменял количество итераций на %d.", s.Config.Iterations))
	case match(q.Data, "/settings/rep"):
		app.bot.EditMessageTextWithKeyboard(q.Message.Chat.ID, q.Message.MessageID, repMenu, repKeyboard)
	case match(q.Data, "/settings/rep/([1-6])", &num):
		// TODO: add symbol to the chosen option
		s.Config.Repeat = num
		app.sessions.Set(q.From.ID, s)
		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Поменял количество повторений на %d.", num))
	case match(q.Data, "/settings/alpha"):
		app.bot.EditMessageTextWithKeyboard(q.Message.Chat.ID, q.Message.MessageID, alphaMenu, alphaKeyboard)
	case match(q.Data, "/settings/alpha/([0-9]+)", &num):
		// TODO: add symbol to the chosen option
		if num < 0 || num > 255 {
			return
		}

		s.Config.Alpha = num
		app.sessions.Set(q.From.ID, s)
		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Поменял значение альфа-канала на %d.", num))
	case match(q.Data, "/settings/alpha/diff"):
		in := make(chan telegram.Message)
		out := make(chan int)

		// Get user input
		s.InputChannel = in
		app.sessions.Set(q.From.ID, s)
		go app.getInputFromUser(
			q.Message.Chat.ID, q.Message.MessageID, 1, 255,
			alphaMenu, alphaKeyboard, in, out,
		)

		s.Config.Alpha = <-out
		s.InputChannel = nil
		app.sessions.Set(q.From.ID, s)
		close(in)

		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Поменял значение альфа-канала на %d.", s.Config.Alpha))
	case match(q.Data, "/settings/ext"):
		app.bot.EditMessageTextWithKeyboard(q.Message.Chat.ID, q.Message.MessageID, extMenu, extKeyboard)
	case match(q.Data, "/settings/ext/(jpg|png|svg|gif)", &slug):
		// TODO: add symbol to the chosen option
		s.Config.Extension = primitive.Extension(slug)
		app.sessions.Set(q.From.ID, s)
		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Поменял расширение на %s.", slug))
	case match(q.Data, "/settings/size"):
		app.bot.EditMessageTextWithKeyboard(q.Message.Chat.ID, q.Message.MessageID, sizeMenu, sizeKeyboard)
	case match(q.Data, "/settings/size/([0-9]+)", &num):
		// TODO: add symbol to the chosen option
		if num < 256 || num > 1920 {
			return
		}

		s.Config.OutputSize = num
		app.sessions.Set(q.From.ID, s)
		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Поменял размеры изображения на %d.", num))
	case match(q.Data, "/settings/size/diff"):
		in := make(chan telegram.Message)
		out := make(chan int)

		// Get user input
		s.InputChannel = in
		app.sessions.Set(q.From.ID, s)
		go app.getInputFromUser(
			q.Message.Chat.ID, q.Message.MessageID, 256, 3840,
			sizeMenu, sizeKeyboard, in, out,
		)

		s.Config.OutputSize = <-out
		s.InputChannel = nil
		app.sessions.Set(q.From.ID, s)
		close(in)

		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Поменял размеры изображения на %d.", s.Config.OutputSize))
	}
}

func (app *application) getInputFromUser(
	chatID, menuMessageID int64,
	min, max int,
	backMenu string,
	backKeyboard telegram.InlineKeyboardMarkup,
	in chan telegram.Message,
	out chan int,
) {
	err := app.bot.EditMessageText(chatID, menuMessageID,
		fmt.Sprintf("Введи число от %d до %d:", min, max))
	if err != nil {
		app.serverError(chatID, err)
		return
	}

	for {
		userMsg := <-in
		app.bot.DeleteMessage(userMsg.Chat.ID, userMsg.MessageID)
		userInput, err := strconv.Atoi(userMsg.Text)
		app.bot.DeleteMessage(userMsg.Chat.ID, userMsg.MessageID)
		if err == nil && userInput >= min && userInput <= max {
			out <- userInput
			close(out)
			err = app.bot.EditMessageTextWithKeyboard(chatID, menuMessageID, backMenu, backKeyboard)
			if err != nil {
				app.serverError(chatID, err)
				return
			}
			return
		}

		err = app.bot.EditMessageText(chatID, menuMessageID,
			fmt.Sprintf("Неверное значение!\nВведи число от %d до %d:", min, max))
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

func (app *application) primitiveWorker(in chan sessions.Session) {
	var s sessions.Session
	for {
		// get next image from queue
		s = <-in

		start := time.Now()
		app.infoLog.Printf("Creating from '%s' for chat '%d': count=%d, mode=%d, alpha=%d, repeat=%d, resolution=%d, extension=%s",
			s.ImgPath, s.ChatID, s.Config.Iterations, s.Config.Shape, s.Config.Alpha, s.Config.Repeat, s.Config.OutputSize, s.Config.Extension)

		// create primitive
		outputPath := fmt.Sprintf("%s/%d_%d.%s", app.outDir, s.ChatID, time.Now().Unix(), s.Config.Extension)
		err := s.Config.Create(s.ImgPath, outputPath)
		if err != nil {
			app.serverError(s.ChatID, err)
			return
		}
		app.infoLog.Printf("Finished creating '%s' for chat '%d'; Output: '%s'; Time: %.1f seconds",
			filepath.Base(s.ImgPath), s.ChatID, filepath.Base(outputPath), time.Since(start).Seconds())

		// send output to the user
		err = app.bot.SendDocument(s.ChatID, outputPath)
		if err != nil {
			app.serverError(s.ChatID, err)
			return
		}
		app.infoLog.Printf("Sent result '%s' to the chat '%d'", filepath.Base(outputPath), s.ChatID)
	}
}
