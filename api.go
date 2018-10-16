package tg

// APIUser represents the "User" JSON structure
type APIUser struct {
	UserID    int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

// ChatType defines the type of chat
type ChatType string

const (
	// ChatTypePrivate is a private chat (between user and bot)
	ChatTypePrivate ChatType = "private"

	// ChatTypeGroup is a group chat (<100 members)
	ChatTypeGroup ChatType = "group"

	// ChatTypeSupergroup is a supergroup chat (>=100 members)
	ChatTypeSupergroup ChatType = "supergroup"

	// ChatTypeChannel is a channel (Read-only)
	ChatTypeChannel ChatType = "channel"
)

// APIChat represents the "Chat" JSON structure
type APIChat struct {
	ChatID    int64    `json:"id"`
	Type      ChatType `json:"type"`
	Title     *string  `json:"title,omitempty"`
	Username  *string  `json:"username,omitempty"`
	FirstName *string  `json:"first_name,omitempty"`
	LastName  *string  `json:"last_name,omitempty"`
}

// APIMessage represents the "Message" JSON structure
type APIMessage struct {
	MessageID         int64          `json:"message_id"`
	User              APIUser        `json:"from"`
	Time              int64          `json:"date"`
	Chat              *APIChat       `json:"chat"`
	FwdUser           *APIUpdate     `json:"forward_from,omitempty"`
	FwdTime           *int           `json:"forward_date,omitempty"`
	ReplyTo           *APIMessage    `json:"reply_to_message,omitempty"`
	Text              *string        `json:"text,omitempty"`
	Audio             *APIAudio      `json:"audio,omitempty"`
	Document          *APIDocument   `json:"document,omitempty"`
	Photo             []APIPhotoSize `json:"photo,omitempty"`
	Sticker           *APISticker    `json:"sticker,omitempty"`
	Video             *APIVideo      `json:"video,omitempty"`
	Voice             *APIVoice      `json:"voice,omitempty"`
	Caption           *string        `json:"caption,omitempty"`
	Contact           *APIContact    `json:"contact,omitempty"`
	Location          *APILocation   `json:"location,omitempty"`
	NewUser           *APIUser       `json:"new_chat_partecipant,omitempty"`
	LeftUser          *APIUser       `json:"left_chat_partecipant,omitempty"`
	PhotoDeleted      *bool          `json:"delete_chat_photo,omitempty"`
	GroupCreated      *bool          `json:"group_chat_created,omitempty"`
	SupergroupCreated *bool          `json:"supergroup_chat_created,omitempty"`
	ChannelCreated    *bool          `json:"channel_chat_created,omitempty"`
	GroupToSuper      *int64         `json:"migrate_to_chat_id,omitempty"`
	GroupFromSuper    *int64         `json:"migrate_from_chat_id,omitempty"`
}

// APIPhotoSize represents the "PhotoSize" JSON structure
type APIPhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize *int   `json:"file_size,omitempty"`
}

// APIAudio represents the "Audio" JSON structure
type APIAudio struct {
	FileID    string  `json:"file_id"`
	Duration  int     `json:"duration"`
	Performer *string `json:"performer,omitempty"`
	Title     *string `json:"title,omitempty"`
	MimeType  *string `json:"mime_type,omitempty"`
	FileSize  *int    `json:"file_size,omitempty"`
}

// APIDocument represents the "Document" JSON structure
type APIDocument struct {
	FileID    string        `json:"file_id"`
	Thumbnail *APIPhotoSize `json:"thumb,omitempty"`
	Filename  string        `json:"file_name"`
	MimeType  *string       `json:"mime_type,omitempty"`
	FileSize  *int          `json:"file_size,omitempty"`
}

// APISticker represents the "Sticker" JSON structure
type APISticker struct {
	FileID    string        `json:"file_id"`
	Width     int           `json:"width"`
	Height    int           `json:"height"`
	Thumbnail *APIPhotoSize `json:"thumb,omitempty"`
	FileSize  *int          `json:"file_size,omitempty"`
}

// APIVideo represents the "Video" JSON structure
type APIVideo struct {
	FileID    string        `json:"file_id"`
	Width     int           `json:"width"`
	Height    int           `json:"height"`
	Duration  int           `json:"duration"`
	Thumbnail *APIPhotoSize `json:"thumb,omitempty"`
	MimeType  *string       `json:"mime_type,omitempty"`
	FileSize  *int          `json:"file_size,omitempty"`
}

// APIVoice represents the "Voice" JSON structure
type APIVoice struct {
	FileID   string  `json:"file_id"`
	Duration int     `json:"duration"`
	MimeType *string `json:"mime_type,omitempty"`
	FileSize *int    `json:"file_size,omitempty"`
}

// APIContact represents the "Contact" JSON structure
type APIContact struct {
	PhoneNumber string  `json:"phone_number"`
	FirstName   string  `json:"first_name"`
	LastName    *string `json:"last_name,omitempty"`
	UserID      *int64  `json:"user_id,omitempty"`
}

// APILocation represents the "Location" JSON structure
type APILocation struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// APIUpdate represents the "Update" JSON structure
type APIUpdate struct {
	UpdateID int64           `json:"update_id"`
	Message  *APIMessage     `json:"message"`
	Inline   *APIInlineQuery `json:"inline_query,omitempty"`
}

// APIFile represents the "File" JSON structure
type APIFile struct {
	FileID string  `json:"file_id"`
	Size   *int    `json:"file_size,omitempty"`
	Path   *string `json:"file_path,omitempty"`
}

// APIResponse represents a response from the Telegram API
type APIResponse struct {
	Ok          bool    `json:"ok"`
	ErrCode     *int    `json:"error_code,omitempty"`
	Description *string `json:"description,omitempty"`
}

// APIInlineQuery represents an inline query from telegram
type APIInlineQuery struct {
	QueryID  string       `json:"id"`
	From     APIUser      `json:"from"`
	Location *APILocation `json:"location,omitempty"`
	Query    string       `json:"query"`
	Offset   string       `json:"offset"`
}

// APIInlineQueryResultPhoto is an image result for an inline query
type APIInlineQueryResultPhoto struct {
	Type        string                   `json:"type"`
	ResultID    string                   `json:"id"`
	PhotoURL    string                   `json:"photo_url"`
	ThumbURL    string                   `json:"thumb_url"`
	Width       int                      `json:"photo_width,omitempty"`
	Height      int                      `json:"photo_height,omitempty"`
	Title       string                   `json:"title,omitempty"`
	Description string                   `json:"description,omitempty"`
	Caption     string                   `json:"caption,omitempty"`
	ParseMode   string                   `json:"parse_mode,omitempty"`
	ReplyMarkup *APIInlineKeyboardMarkup `json:"reply_markup,omitempty"`
	//TODO inputMessageContent
}

type APIInlineKeyboardMarkup struct {
	InlineKeyboard interface{} `json:"inline_keyboard"`
}

// APIInlineKeyboardButton is an inline message button
type APIInlineKeyboardButton struct {
	Text string `json:"text"`
	URL  string `json:"url,omitempty"`
}

// APIInputMediaPhoto is a media photo element (already on telegram servers or via HTTP URL) for albums and other cached pictures
type APIInputMediaPhoto struct {
	Type      string `json:"type"`
	Media     string `json:"media"`
	Caption   string `json:"caption,omitempty"`
	ParseMode string `json:"parse_mode,omitempty"`
}
