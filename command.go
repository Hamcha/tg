package tg

// BrokerUpdateType distinguishes update types coming from the broker
type BrokerUpdateType string

const (
	// BMessage is a message update (mostly webhook updates)
	BMessage BrokerUpdateType = "message"

	// BFile is a file retrieval response update
	BFile BrokerUpdateType = "file"

	// BError is an error the broker occurred while fulfilling a request
	BError BrokerUpdateType = "error"
)

// BrokerUpdate is what is sent by the broker as update
type BrokerUpdate struct {
	Type     BrokerUpdateType
	Callback *int       `json:",omitempty"`
	Error    *string    `json:",omitempty"`
	Data     *APIUpdate `json:",omitempty"`
	Bytes    *string    `json:",omitempty"`
}

// ClientCommandType distinguishes requests sent by clients to the broker
type ClientCommandType string

const (
	// CmdSendTextMessage requests the broker to send a text message to a chat
	CmdSendTextMessage ClientCommandType = "sendText"

	// CmdSendPhoto requests the broker to send a photo to a chat
	CmdSendPhoto ClientCommandType = "sendPhoto"

	// CmdForwardMessage requests the broker to forward a message between chats
	CmdForwardMessage ClientCommandType = "forwardMessage"

	// CmdGetFile requests the broker to get a file from Telegram
	CmdGetFile ClientCommandType = "getFile"

	// CmdSendChatAction requests the broker to set a chat action (typing, etc.)
	CmdSendChatAction ClientCommandType = "sendChatAction"

	// CmdAnswerInlineQuery requests the broker sends results of an inline query
	CmdAnswerInlineQuery ClientCommandType = "answerInlineQuery"
)

// ClientTextMessageData is the required data for a CmdSendTextMessage request
type ClientTextMessageData struct {
	ChatID  int64
	Text    string
	ReplyID *int64 `json:",omitempty"`
}

// ClientPhotoData is the required data for a CmdSendPhoto request
type ClientPhotoData struct {
	ChatID   int64
	Bytes    string
	Filename string
	Caption  string `json:",omitempty"`
	ReplyID  *int64 `json:",omitempty"`
}

// ClientForwardMessageData is the required data for a CmdForwardMessage request
type ClientForwardMessageData struct {
	ChatID     int64
	FromChatID int64
	MessageID  int64
}

// ClientChatActionData is the required data for a CmdSendChatAction request
type ClientChatActionData struct {
	ChatID int64
	Action ChatAction
}

// ChatAction is the action name for CmdSendChatAction requests
type ChatAction string

const (
	ActionTyping            ChatAction = "typing"
	ActionUploadingPhoto    ChatAction = "upload_photo"
	ActionRecordingVideo    ChatAction = "record_video"
	ActionUploadingVideo    ChatAction = "upload_video"
	ActionRecordingAudio    ChatAction = "record_audio"
	ActionUploadingAudio    ChatAction = "upload_audio"
	ActionUploadingDocument ChatAction = "upload_document"
	ActionFindingLocation   ChatAction = "find_location"
)

// FileRequestData is the required data for a CmdGetFile request
type FileRequestData struct {
	FileID string
}

// ClientCommand is a request sent by clients to the broker
type ClientCommand struct {
	Type               ClientCommandType
	TextMessageData    *ClientTextMessageData    `json:",omitempty"`
	PhotoData          *ClientPhotoData          `json:",omitempty"`
	ForwardMessageData *ClientForwardMessageData `json:",omitempty"`
	ChatActionData     *ClientChatActionData     `json:",omitempty"`
	InlineQueryResults *InlineQueryResponse      `json:",omitempty"`
	FileRequestData    *FileRequestData          `json:",omitempty"`
	Callback           *int                      `json:",omitempty"`
}

// InlineQueryResponse is the response to an inline query
type InlineQueryResponse struct {
	QueryID    string
	Results    []interface{}
	CacheTime  *int   `json:",omitempty"`
	IsPersonal bool   `json:",omitempty"`
	NextOffset string `json:",omitempty"`
	PMText     string `json:",omitempty"`
	PMParam    string `json:",omitempty"`
}
