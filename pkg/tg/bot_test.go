package tg

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

var chatID int64

func getBot(t *testing.T) *Bot {
	token := os.Getenv("TEST_BOT_TOKEN")
	if token == "" {
		t.Fatal("You must provide token for the test bot via TEST_BOT_TOKEN environmental variable!")
	}

	id := os.Getenv("TEST_CHAT_ID")
	if id == "" {
		t.Fatal("You must provide receiver chat ID via TEST_CHAT_ID environmental variable!")
	}

	num, err := strconv.Atoi(id)
	if err != nil {
		t.Fatalf("Incorrect value in TEST_CHAT_ID: %v", err)
	}
	chatID = int64(num)

	return &Bot{Token: token}
}

func TestBot_GetUpdates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	bot := getBot(t)

	_, err := bot.GetUpdates(0, 100, 0, []string{})
	if err != nil {
		t.Errorf("Error getting updates: %v", err)
	}
}

func TestBot_AnswerCallbackQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

func TestBot_SendMessageWhenOnlyText(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	bot := getBot(t)

	_, err := bot.SendMessage(chatID, "Test message")
	if err != nil {
		t.Errorf("Error sending the message: %v", err)
	}
}

func TestBot_SendMessageWhenWithKeyboard(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	bot := getBot(t)

	_, err := bot.SendMessage(chatID, "Test message", InlineKeyboardMarkup{
		InlineKeyboard: [][]InlineKeyboardButton{{{Text: "Hello", CallbackData: "hello"}}},
	})
	if err != nil {
		t.Errorf("Error sending the message: %v", err)
	}
}

func TestBot_SendDocument(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	bot := getBot(t)

	path := "test_file.txt"
	err := os.WriteFile(path, []byte("Hello World"), 0600)
	if err != nil {
		t.Fatalf("Error while creating test file: %v", err)
	}
	defer os.Remove(path)

	err = bot.SendDocument(chatID, path)
	if err != nil {
		t.Errorf("Error sending document: %v", err)
	}
}

func TestBot_SendDocumentWhenPathIsIncorrect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	bot := getBot(t)

	path := "test_file.txt"

	err := bot.SendDocument(chatID, path)
	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Error sending non-existing document: %v", err)
	}
}

func TestBot_EditMessageText(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	bot := getBot(t)

	// send message
	msg, err := bot.SendMessage(chatID, "Test message that should be edited.")
	if err != nil {
		t.Fatalf("Error sending message: %v", err)
	}

	// edit message
	err = bot.EditMessageText(chatID, msg.MessageID, "Edited test message.")
	if err != nil {
		t.Fatalf("Error editing message: %v", err)
	}
}

func TestBot_DeleteMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	bot := getBot(t)

	// send message
	msg, err := bot.SendMessage(chatID, "Test message that should be deleted.")
	if err != nil {
		t.Fatalf("Error sending message: %v", err)
	}

	// delete message
	err = bot.DeleteMessage(chatID, msg.MessageID)
	if err != nil {
		t.Fatalf("Error deleting message: %v", err)
	}
}
