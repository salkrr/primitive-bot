package tg

// InlineKeyboardButton object represents one button of an inline keyboard.
type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

// InlineKeyboardMarkup object represents an inline keyboard that appears right next to the message it belongs to.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// Chat object represents a chat.
type Chat struct {
	ID int64 `json:"id"`
}

// User object represents a Telegram user or bot.
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// File object represents a file ready to be downloaded.
type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FilePath     string `json:"file_path"`
}

// Document object represents a general file (as opposed to photos, voice messages and audio files).
type Document struct {
	FileID   string `json:"file_id"`
	UniqueID string `json:"file_unique_id"`
	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
	FileSize int    `json:"file_size"`
}

// PhotoSize object represents one size of a photo or a file / sticker thumbnail.
type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size"`
}

// Message object represents a message.
type Message struct {
	MessageID int64       `json:"message_id"`
	From      User        `json:"from"`
	Chat      Chat        `json:"chat"`
	Date      int64       `json:"date"`
	Text      string      `json:"text"`
	Photo     []PhotoSize `json:"photo"`
	Document  Document    `json:"document"`
}

// CallbackQuery object represents an incoming callback query from a callback button in an inline keyboard.
type CallbackQuery struct {
	ID      string  `json:"id"`
	From    User    `json:"from"`
	Message Message `json:"message"`
	Data    string  `json:"data"`
}

// Update object represents an incoming update.
type Update struct {
	UpdateID      int64         `json:"update_id"`
	Message       Message       `json:"message"`
	CallbackQuery CallbackQuery `json:"callback_query"`
}

// Response is a response to a request to the Telegram API.
type Response struct {
	Ok          bool        `json:"ok"`
	ErrorCode   int         `json:"error_code"`
	Description string      `json:"description"`
	Result      interface{} `json:"result"`
}
