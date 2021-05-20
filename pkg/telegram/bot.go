// Package telegram implements some of the Telegram API methods
package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	baseBotURL  = "https://api.telegram.org/bot"
	baseFileURL = "https://api.telegram.org/file/bot"
)

type Bot struct {
	Token string
}

func (b *Bot) GetUpdates(offset int64) ([]Update, error) {
	u, err := url.Parse(fmt.Sprintf("%s%s/getUpdates", baseBotURL, b.Token))
	if err != nil {
		return nil, err
	}

	reqJSON, err := json.Marshal(map[string]interface{}{
		"offset":          offset,
		"limit":           100,
		"timeout":         20,
		"allowed_updates": []string{"message", "callback_query"},
	})
	if err != nil {
		return nil, err
	}

	resp, err := b.sendPostRequest(u.String(), "application/json", bytes.NewBuffer(reqJSON))
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

func (b *Bot) AnswerCallbackQuery(callbackID, text string) error {
	u, err := url.Parse(fmt.Sprintf("%s%s/answerCallbackQuery", baseBotURL, b.Token))
	if err != nil {
		return err
	}

	q := url.Values{}
	q.Set("callback_query_id", callbackID)
	q.Set("text", text)
	u.RawQuery = q.Encode()

	_, err = b.sendGetRequest(u.String())
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) EditMessageText(chatID, messageID int64, text string, keyboard ...InlineKeyboardMarkup) error {
	u, err := url.Parse(fmt.Sprintf("%s%s/editMessageText", baseBotURL, b.Token))
	if err != nil {
		return err
	}

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

	_, err = b.sendPostRequest(u.String(), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	return nil
}

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

	u, err := url.Parse(fmt.Sprintf("%s%s/sendMessage", baseBotURL, b.Token))
	if err != nil {
		return Message{}, err
	}

	resp, err := b.sendPostRequest(u.String(), "application/json", bytes.NewBuffer(jsonBody))
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

func (b *Bot) SendPhoto(chatID int64, photoPath string) error {
	w, formBody, err := createMultipartForm("photo", photoPath)
	if err != nil {
		return err
	}

	u, err := url.Parse(fmt.Sprintf("%s%s/sendPhoto", baseBotURL, b.Token))
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

	var respContent Response
	err = json.Unmarshal(body, &respContent)
	if err != nil {
		return err
	}

	if !respContent.Ok {
		return fmt.Errorf("error code: %v; description: %s", respContent.ErrorCode, respContent.Description)
	}

	return nil
}

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

	var respContent Response
	err = json.Unmarshal(body, &respContent)
	if err != nil {
		return err
	}

	if !respContent.Ok {
		return fmt.Errorf("error code: %v; description: %s", respContent.ErrorCode, respContent.Description)
	}

	return nil
}

func (b *Bot) DeleteMessage(chatID, messageID int64) error {
	u, err := url.Parse(fmt.Sprintf("%s%s/deleteMessage", baseBotURL, b.Token))
	if err != nil {
		return err
	}

	q := url.Values{}
	q.Set("chat_id", fmt.Sprint(chatID))
	q.Set("message_id", fmt.Sprint(messageID))
	u.RawQuery = q.Encode()

	_, err = b.sendGetRequest(u.String())
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) GetFile(fileID string) (File, error) {
	u, err := url.Parse(fmt.Sprintf("%s%s/getFile", baseBotURL, b.Token))
	if err != nil {
		return File{}, err
	}

	q := url.Values{}
	q.Set("file_id", fileID)
	u.RawQuery = q.Encode()

	resp, err := b.sendGetRequest(u.String())
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

func (s *Bot) DownloadFile(fileID string) ([]byte, error) {
	file, err := s.GetFile(fileID)
	if err != nil {
		return nil, err
	}

	if file.FileID == "" {
		return nil, fmt.Errorf("file doesn't have an ID")
	}

	u, err := url.Parse(fmt.Sprintf("%s%s/%s", baseFileURL, s.Token, file.FilePath))
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

func (b *Bot) sendPostRequest(url string, contentType string, formBody *bytes.Buffer) (Response, error) {
	resp, err := http.Post(url, contentType, formBody)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var respContent Response
	err = json.Unmarshal(body, &respContent)
	if err != nil {
		return Response{}, err
	}

	if !respContent.Ok {
		return Response{}, fmt.Errorf("error code: %v; description: %s", respContent.ErrorCode, respContent.Description)
	}

	return respContent, nil
}

func (b *Bot) sendGetRequest(url string) (Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var respContent Response
	err = json.Unmarshal(body, &respContent)
	if err != nil {
		return Response{}, err
	}

	if !respContent.Ok {
		return Response{}, fmt.Errorf("error code: %v; description: %s", respContent.ErrorCode, respContent.Description)
	}

	return respContent, nil
}
