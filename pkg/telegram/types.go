package telegram

type KeyboardButton struct {
	Text string `json:"text"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`
	OneTimeKeyboard bool               `json:"one_time_keyboard"`
	Selective       bool               `json:"selective"`
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FilePath     string `json:"file_path"`
}

type Document struct {
	FileID   string `json:"file_id"`
	UniqueID string `json:"file_unique_id"`
	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
	FileSize int    `json:"file_size"`
}

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size"`
}

type Message struct {
	MessageID int64       `json:"message_id"`
	From      User        `json:"from"`
	Chat      Chat        `json:"chat"`
	Date      int64       `json:"date"`
	Text      string      `json:"text"`
	Photo     []PhotoSize `json:"photo"`
	Document  Document    `json:"document"`
}

type Update struct {
	UpdateID int64   `json:"update_id"`
	Message  Message `json:"message"`
}

type Response struct {
	Ok          bool        `json:"ok"`
	ErrorCode   int         `json:"error_code"`
	Description string      `json:"description"`
	Result      interface{} `json:"result"`
}
