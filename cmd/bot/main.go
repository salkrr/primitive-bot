package main

import (
	"flag"
	"log"
	"os"

	"github.com/lazy-void/primitive-bot/pkg/sessions"
	"github.com/lazy-void/primitive-bot/pkg/telegram"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	bot      *telegram.Bot
	sessions *sessions.ActiveSessions
}

func main() {
	token := flag.String("token", "", "The token for the Telegram Bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("You need to provide the token for the Telegram Bot!")
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := application{
		infoLog:  infoLog,
		errorLog: errorLog,
		bot:      &telegram.Bot{Token: *token},
		sessions: sessions.New(),
	}

	infoLog.Printf("Starting to listen for updates...")
	app.listenAndServe()
}

func (app *application) listenAndServe() {
	// create worker
	ch := make(chan sessions.Session, 100)
	go app.primitiveWorker(ch)

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
				go app.handleMessage(u.Message, ch)
			} else if u.CallbackQuery.ID != "" {
				app.infoLog.Printf("Got callback query with data '%s' from the user '%s' with ID '%d'",
					u.CallbackQuery.Data, u.CallbackQuery.From.FirstName, u.CallbackQuery.From.ID)
				go app.handleCallbackQuery(u.CallbackQuery, ch)
			}
		}

		offset = updates[numUpdates-1].UpdateID + 1
	}
}
