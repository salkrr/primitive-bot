package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/lazy-void/primitive-bot/pkg/queue"
	"github.com/lazy-void/primitive-bot/pkg/sessions"
	"github.com/lazy-void/primitive-bot/pkg/telegram"
)

type application struct {
	infoLog         *log.Logger
	errorLog        *log.Logger
	inDir           string
	outDir          string
	operationsLimit int
	bot             *telegram.Bot
	sessions        *sessions.ActiveSessions
	queue           *queue.Queue
}

func main() {
	token := flag.String("token", "", "The token for the Telegram Bot")
	inDir := flag.String("i", "inputs", "Name of the directory where inputs will be stored")
	outDir := flag.String("o", "outputs", "Name of the directory where outputs will be stored")
	operationsLimit := flag.Int("limit", 5, "The number of operations that the user can add to the queue.")
	flag.Parse()

	if *token == "" {
		log.Fatal("You need to provide the token for the Telegram Bot!")
	}

	if err := os.MkdirAll(*inDir, 0664); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(*outDir, 0664); err != nil {
		log.Fatal(err)
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := application{
		infoLog:         infoLog,
		errorLog:        errorLog,
		inDir:           *inDir,
		outDir:          *outDir,
		operationsLimit: *operationsLimit,
		bot:             &telegram.Bot{Token: *token},
		sessions:        sessions.NewActiveSessions(),
		queue:           queue.New(),
	}

	infoLog.Printf("Starting to listen for updates...")
	app.listenAndServe()
}

func (app *application) listenAndServe() {
	go app.worker()

	offset := int64(0)
	for {
		updates, err := app.bot.GetUpdates(offset)
		if err != nil {
			app.errorLog.Print(err)
			continue
		}

		numUpdates := len(updates)
		if numUpdates == 0 {
			continue
		}

		for _, u := range updates {
			if u.Message.MessageID > 0 {
				app.infoLog.Printf("Got message with text '%s' from the user '%s' with ID '%d'",
					u.Message.Text, u.Message.From.FirstName, u.Message.From.ID)
				go app.handleMessage(u.Message)
				continue
			}

			app.infoLog.Printf("Got callback query with data '%s' from the user '%s' with ID '%d'",
				u.CallbackQuery.Data, u.CallbackQuery.From.FirstName, u.CallbackQuery.From.ID)
			go app.handleCallbackQuery(u.CallbackQuery)
		}

		offset = updates[numUpdates-1].UpdateID + 1
	}
}

func (app *application) worker() {
	for {
		// get next operation
		op, ok := app.queue.Peek()
		if !ok {
			time.Sleep(1 * time.Second)
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
