package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lazy-void/primitive-bot/pkg/primitive"
	"github.com/lazy-void/primitive-bot/pkg/sessions"
	"github.com/lazy-void/primitive-bot/pkg/telegram"
)

func (app *application) handleMessage(m telegram.Message, out chan sessions.Session) {
	if m.Photo == nil {
		err := app.bot.SendMessage(m.Chat.ID, helpMessage)
		if err != nil {
			app.serverError(m.Chat.ID, err)
		}
		return
	}

	// Download image
	file := m.Photo[2]
	img, err := app.bot.DownloadFile(file.FileID)
	if err != nil {
		app.serverError(m.Chat.ID, fmt.Errorf("couldn't download image: %s", err))
		return
	}
	path := fmt.Sprintf("./inputs/%s.jpg", file.FileUniqueID)
	os.WriteFile(path, img, 0644)

	// Create session
	app.sessions.Set(m.From.ID, sessions.Session{
		ChatID:       m.Chat.ID,
		ImgMessageID: m.MessageID,
		ImgPath:      path,
		Config:       primitive.NewConfig(),
	})
	app.bot.SendMessageWithInlineKeyboard(m.Chat.ID, rootMenu, rootKeyboard)
}

func (app *application) handleCallbackQuery(q telegram.CallbackQuery, out chan sessions.Session) {
	defer app.bot.AnswerCallbackQuery(q.ID, "")

	s, err := app.sessions.Get(q.From.ID)
	if err == sessions.ErrNoSession {
		return
	} else if err != nil {
		app.serverError(q.Message.Chat.ID, err)
	}

	var num int
	var slug string
	switch {
	case match(q.Data, "/"):
		app.bot.EditMessageText(q.Message.Chat.ID, q.Message.MessageID, rootMenu, rootKeyboard)

	case match(q.Data, "/start"):
		out <- s
		app.sessions.Delete(q.From.ID)
		app.bot.DeleteMessage(q.Message.Chat.ID, q.Message.MessageID)
		app.bot.SendMessage(q.Message.Chat.ID, fmt.Sprintf(enqueuedMessage,
			strings.ToLower(shapeNames[s.Config.Shape]), s.Config.Iterations,
			s.Config.Repeat, s.Config.Alpha, s.Config.Extension, s.Config.OutputSize))

	case match(q.Data, "/settings/shape"):
		app.bot.EditMessageText(q.Message.Chat.ID, q.Message.MessageID, shapesMenu, shapesKeyboard)

	case match(q.Data, "/settings/shape/([0-8])", &num):
		// TODO: add symbol to the chosen option
		s.Config.Shape = primitive.Shape(num)
		app.sessions.Set(q.From.ID, s)
		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Будут использоваться фигуры: %s.", strings.ToLower(shapeNames[primitive.Shape(num)])))

	case match(q.Data, "/settings/iter"):
		app.bot.EditMessageText(q.Message.Chat.ID, q.Message.MessageID, iterMenu, iterKeyboard)

	case match(q.Data, "/settings/iter/([0-9]+)", &num):
		// TODO: add symbol to the chosen option
		s.Config.Iterations = num
		app.sessions.Set(q.From.ID, s)
		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Поменял количество итераций на %d.", num))

	case match(q.Data, "/settings/rep"):
		app.bot.EditMessageText(q.Message.Chat.ID, q.Message.MessageID, repMenu, repKeyboard)

	case match(q.Data, "/settings/rep/([1-6])", &num):
		// TODO: add symbol to the chosen option
		s.Config.Repeat = num
		app.sessions.Set(q.From.ID, s)
		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Поменял количество повторений на %d.", num))

	case match(q.Data, "/settings/alpha"):
		app.bot.EditMessageText(q.Message.Chat.ID, q.Message.MessageID, alphaMenu, alphaKeyboard)

	case match(q.Data, "/settings/alpha/([0-9]+)", &num):
		// TODO: add symbol to the chosen option
		if num < 0 || num > 255 {
			return
		}

		s.Config.Alpha = num
		app.sessions.Set(q.From.ID, s)
		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Поменял значение альфа-канала на %d.", num))

	case match(q.Data, "/settings/ext"):
		app.bot.EditMessageText(q.Message.Chat.ID, q.Message.MessageID, extMenu, extKeyboard)

	case match(q.Data, "/settings/ext/(jpg|png|svg|gif)", &slug):
		// TODO: add symbol to the chosen option
		s.Config.Extension = primitive.Extension(slug)
		app.sessions.Set(q.From.ID, s)
		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Поменял расширение на %s.", slug))

	case match(q.Data, "/settings/size"):
		app.bot.EditMessageText(q.Message.Chat.ID, q.Message.MessageID, sizeMenu, sizeKeyboard)

	case match(q.Data, "/settings/size/([0-9]+)", &num):
		// TODO: add symbol to the chosen option
		if num < 0 || num > 1920 {
			return
		}

		s.Config.OutputSize = num
		app.sessions.Set(q.From.ID, s)
		app.bot.AnswerCallbackQuery(q.ID, fmt.Sprintf("Поменял размеры изображения на %d.", num))
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
		outputPath := fmt.Sprintf("./outputs/%d_%d.%s", s.ChatID, time.Now().Unix(), s.Config.Extension)
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
