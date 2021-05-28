// Package tg implements some of the Telegram API methods
package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	baseBotURL  = "https://api.telegram.org/bot"
	baseFileURL = "https://api.telegram.org/file/bot"

	urlencodedContentType = "application/x-www-form-urlencoded"
	jsonContentType       = "application/json"
)

// Bot is an instance of a Telegram bot.
type Bot struct {
	Token string
}

// GetUpdates implements Telegram's getUpdates method.
func (b *Bot) GetUpdates(offset int64, limit, timeout int, allowedUpdates []string) ([]Update, error) {
	reqJSON, err := json.Marshal(map[string]interface{}{
		"offset":          offset,
		"limit":           limit,
		"timeout":         timeout,
		"allowed_updates": allowedUpdates,
	})
	if err != nil {
		return nil, err
	}

	resp, err := b.makeRequest("/getUpdates", jsonContentType, bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}

	if resp.Result == nil {
		return []Update{}, nil
	}

	resultJSON, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}

	var result []Update
	err = json.Unmarshal(resultJSON, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// AnswerCallbackQuery implements Telegram's answerCallbackQuery method.
func (b *Bot) AnswerCallbackQuery(callbackID, text string) error {
	q := url.Values{}
	q.Set("callback_query_id", callbackID)
	q.Set("text", text)
	_, err := b.makeRequest("/answerCallbackQuery", urlencodedContentType, strings.NewReader(q.Encode()))
	if err != nil {
		return err
	}

	return nil
}

// EditMessageText implements Telegram's editMessageText method.
func (b *Bot) EditMessageText(chatID, messageID int64, text string, keyboard ...InlineKeyboardMarkup) error {
	params := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
		"text":       text,
	}

	if len(keyboard) > 0 {
		params["reply_markup"] = keyboard[0]
	}

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = b.makeRequest("/editMessageText", jsonContentType, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	return nil
}

// SendMessage implements Telegram's sendMessage method.
func (b *Bot) SendMessage(chatID int64, message string, keyboard ...InlineKeyboardMarkup) (Message, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
		"text":    message,
	}
	if len(keyboard) > 0 {
		params["reply_markup"] = keyboard[0]
	}

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return Message{}, err
	}

	resp, err := b.makeRequest("/sendMessage", jsonContentType, bytes.NewBuffer(jsonBody))
	if err != nil {
		return Message{}, err
	}

	resultJSON, err := json.Marshal(resp.Result)
	if err != nil {
		return Message{}, err
	}

	var result Message
	err = json.Unmarshal(resultJSON, &result)
	if err != nil {
		return Message{}, err
	}

	return result, nil
}

// SendDocument implements Telegram's sendDocument method.
func (b *Bot) SendDocument(chatID int64, documentPath string) error {
	w, formBody, err := createMultipartForm("document", documentPath)
	if err != nil {
		return err
	}

	u, err := url.Parse(fmt.Sprintf("%s%s/sendDocument", baseBotURL, b.Token))
	if err != nil {
		return err
	}

	q := url.Values{}
	q.Set("chat_id", fmt.Sprint(chatID))
	u.RawQuery = q.Encode()

	resp, err := http.Post(u.String(), w.FormDataContentType(), formBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var respContent APIResponse
	err = json.Unmarshal(body, &respContent)
	if err != nil {
		return err
	}

	if !respContent.Ok {
		return fmt.Errorf("error code: %v; description: %s", respContent.ErrorCode, respContent.Description)
	}

	return nil
}

// DeleteMessage implements Telegram's deleteMessage method.
func (b *Bot) DeleteMessage(chatID, messageID int64) error {
	q := url.Values{}
	q.Set("chat_id", fmt.Sprint(chatID))
	q.Set("message_id", fmt.Sprint(messageID))

	_, err := b.makeRequest("/deleteMessage", urlencodedContentType, strings.NewReader(q.Encode()))
	if err != nil {
		return err
	}

	return nil
}

// GetFile implements Telegram's getFile method.
func (b *Bot) GetFile(fileID string) (File, error) {
	q := url.Values{}
	q.Set("file_id", fileID)

	resp, err := b.makeRequest("/getFile", urlencodedContentType, strings.NewReader(q.Encode()))
	if err != nil {
		return File{}, err
	}

	resultJSON, err := json.Marshal(resp.Result)
	if err != nil {
		return File{}, err
	}

	var result File
	err = json.Unmarshal(resultJSON, &result)
	if err != nil {
		return File{}, err
	}

	return result, nil
}

// DownloadFile downloads file from the Telegram server.
func (b *Bot) DownloadFile(fileID string) ([]byte, error) {
	file, err := b.GetFile(fileID)
	if err != nil {
		return nil, err
	}

	if file.FileID == "" {
		return nil, fmt.Errorf("file doesn't have an ID")
	}

	u, err := url.Parse(fmt.Sprintf("%s%s/%s", baseFileURL, b.Token, file.FilePath))
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (b *Bot) makeRequest(method string, contentType string, body io.Reader) (APIResponse, error) {
	u := fmt.Sprint(baseBotURL + b.Token + method)

	resp, err := http.Post(u, contentType, body) // #nosec
	if err != nil {
		return APIResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, err
	}

	var respContent APIResponse
	err = json.Unmarshal(respBody, &respContent)
	if err != nil {
		return APIResponse{}, err
	}

	if !respContent.Ok {
		return APIResponse{}, fmt.Errorf("error code: %v; description: %s", respContent.ErrorCode, respContent.Description)
	}

	return respContent, nil
}
