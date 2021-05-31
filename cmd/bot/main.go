package main

//go:generate gotext update -out=catalog.go -lang=en,ru

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/lazy-void/primitive-bot/pkg/menu"

	"golang.org/x/text/language"

	"golang.org/x/text/message"

	"github.com/lazy-void/primitive-bot/pkg/primitive"

	"github.com/lazy-void/primitive-bot/pkg/queue"
	"github.com/lazy-void/primitive-bot/pkg/sessions"
	"github.com/lazy-void/primitive-bot/pkg/tg"
)

var (
	token           string
	inDir           string
	outDir          string
	logPath         string
	operationsLimit int
	maxIter         int
	maxSize         int
	workers         int
	timeout         time.Duration
	lang            language.Tag
)

type application struct {
	infoLog         *log.Logger
	errorLog        *log.Logger
	printer         *message.Printer
	inDir           string
	outDir          string
	operationsLimit int
	maxIter         int
	maxSize         int
	workers         int
	bot             *tg.Bot
	sessions        *sessions.ActiveSessions
	queue           *queue.Queue
}

func init() {
	flag.StringVar(&token, "token", "", "The token for the Telegram Bot.")
	flag.StringVar(&inDir, "i", "inputs", "Name of the directory where inputs will be stored.")
	flag.StringVar(&outDir, "o", "outputs", "Name of the directory where outputs will be stored.")
	flag.StringVar(&logPath, "log", "", "Path to the previous log file. It is used to restore queue.")
	flag.IntVar(&workers, "w", runtime.NumCPU(), "Numbers of parallel workers used to create primitive image.")
	flag.IntVar(&operationsLimit, "limit", 5, "The number of operations that the user can add to the queue.")
	flag.IntVar(&maxIter, "iter", 2000, "Maximum iterations that user can specify.")
	flag.IntVar(&maxSize, "size", 3840, "Maximum image size that user can specify.")
	flag.DurationVar(&timeout, "timeout", 30*time.Minute,
		"The number of minutes that a session can be inactive before it's terminated.")
	flag.Func("lang", `Language of the bot (en, ru). (default "en")`, func(s string) error {
		if s != "en" && s != "ru" {
			return errors.New("incorrect language")
		}

		lang = language.MustParse(s)
		return nil
	})
}

func main() {
	flag.Parse()

	if token == "" {
		log.Fatal("You need to provide the token for the Telegram Bot!")
	}
	if lang.String() == "und" {
		lang = language.MustParse("en")
	}

	// restore queue if needed
	q := queue.New()
	if logPath != "" {
		err := restoreQueue(logPath, q, workers)
		if err != nil {
			log.Fatalf("Error restoring queue from the log: %v", err)
		}
	}

	// create directories for input and output
	if err := os.MkdirAll(inDir, 0664); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(outDir, 0664); err != nil {
		log.Fatal(err)
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// initialize localization
	printer := message.NewPrinter(lang)
	menu.InitText(printer)

	app := application{
		infoLog:         infoLog,
		errorLog:        errorLog,
		printer:         printer,
		inDir:           inDir,
		outDir:          outDir,
		operationsLimit: operationsLimit,
		maxIter:         maxIter,
		maxSize:         maxSize,
		workers:         workers,
		bot:             &tg.Bot{Token: token},
		sessions:        sessions.NewActiveSessions(timeout, 5*time.Minute),
		queue:           q,
	}

	infoLog.Printf("Starting to listen for the updates...")
	app.listenAndServe()
}

func restoreQueue(logPath string, q *queue.Queue, workers int) (err error) {
	f, err := os.Open(filepath.Clean(logPath))
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()

	regex := regexp.MustCompile(`INFO\t\b.*?\b \b.*?\b (.*)`)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Remove log prefix
		matches := regex.FindStringSubmatch(scanner.Text())
		if len(matches) < 2 {
			continue
		}

		msg := matches[1]
		if strings.HasPrefix(msg, "Enqueued:") {
			op := queue.Operation{Config: primitive.New(workers)}
			_, err := fmt.Sscanf(msg, enqueuedLogMessage, &op.UserID, &op.ImgPath,
				&op.Config.Iterations, &op.Config.Shape, &op.Config.Alpha,
				&op.Config.Repeat, &op.Config.OutputSize, &op.Config.Extension)
			if err != nil {
				return err
			}

			q.Enqueue(op)
			continue
		}

		if strings.HasPrefix(msg, "Sent:") {
			q.Dequeue()
		}
	}

	return scanner.Err()
}
