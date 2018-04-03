package tg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// Broker is a broker connection handler with callback management functions
type Broker struct {
	Socket    net.Conn
	Callbacks []BrokerCallback

	cbFree int
}

// ConnectToBroker creates a Broker connection
func ConnectToBroker(brokerAddr string) (*Broker, error) {
	sock, err := net.Dial("tcp", brokerAddr)
	if err != nil {
		return nil, err
	}

	broker := new(Broker)
	broker.Socket = sock
	broker.Callbacks = make([]BrokerCallback, 0)
	broker.cbFree = 0
	return broker, nil
}

// Close closes a broker connection
func (b *Broker) Close() {
	b.Socket.Close()
}

// SendTextMessage sends a HTML-styles text message to a chat.
// A reply_to message ID can be specified as optional parameter.
func (b *Broker) SendTextMessage(chat *APIChat, text string, original *int64) {
	b.sendCmd(ClientCommand{
		Type: CmdSendTextMessage,
		TextMessageData: &ClientTextMessageData{
			Text:    text,
			ChatID:  chat.ChatID,
			ReplyID: original,
		},
	})
}

// SendPhoto sends a photo with an optional caption to a chat.
// A reply_to message ID can be specified as optional parameter.
func (b *Broker) SendPhoto(chat *APIChat, data []byte, filename string, caption string, original *int64) {
	b.sendCmd(ClientCommand{
		Type: CmdSendPhoto,
		PhotoData: &ClientPhotoData{
			ChatID:   chat.ChatID,
			Filename: filename,
			Bytes:    base64.StdEncoding.EncodeToString(data),
			Caption:  caption,
			ReplyID:  original,
		},
	})
}

// ForwardMessage forwards a message between chats.
func (b *Broker) ForwardMessage(chat *APIChat, message APIMessage) {
	b.sendCmd(ClientCommand{
		Type: CmdForwardMessage,
		ForwardMessageData: &ClientForwardMessageData{
			ChatID:     chat.ChatID,
			FromChatID: message.Chat.ChatID,
			MessageID:  message.MessageID,
		},
	})
}

// SendChatAction sets a chat action for 5 seconds or less (canceled at first message sent)
func (b *Broker) SendChatAction(chat *APIChat, action ChatAction) {
	b.sendCmd(ClientCommand{
		Type: CmdSendChatAction,
		ChatActionData: &ClientChatActionData{
			ChatID: chat.ChatID,
			Action: action,
		},
	})
}

// GetFile sends a file retrieval request to the Broker.
// This function is asynchronous as data will be delivered to the given callback.
func (b *Broker) GetFile(fileID string, fn BrokerCallback) int {
	cid := b.RegisterCallback(fn)
	b.sendCmd(ClientCommand{
		Type: CmdGetFile,
		FileRequestData: &FileRequestData{
			FileID: fileID,
		},
		Callback: &cid,
	})
	return cid
}

// RegisterCallback assigns a callback ID to the given callback and puts it on the callback list.
// This function should never be called by clients.
func (b *Broker) RegisterCallback(fn BrokerCallback) int {
	cblen := len(b.Callbacks)
	// List is full, append to end
	if b.cbFree == cblen {
		b.Callbacks = append(b.Callbacks, fn)
		b.cbFree++
		return cblen
	}
	// List is not full, use empty slot and find next one
	id := b.cbFree
	b.Callbacks[id] = fn
	next := b.cbFree + 1
	for ; next < cblen; next++ {
		if b.Callbacks[next] == nil {
			break
		}
	}
	b.cbFree = next
	return id
}

// RemoveCallback removes a callback from the callback list by ID.
// This function should never be called by clients.
func (b *Broker) RemoveCallback(id int) {
	b.Callbacks[id] = nil
	if id < b.cbFree {
		b.cbFree = id
	}
	b.resizeCbArray()
}

// SpliceCallback retrieves a callback by ID and removes it from the list.
// This function should never be called by clients.
func (b *Broker) SpliceCallback(id int) BrokerCallback {
	defer b.RemoveCallback(id)
	return b.Callbacks[id]
}

func (b *Broker) sendCmd(cmd ClientCommand) {
	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("[sendCmd] JSON Encode error: %s\n", err.Error())
	}
	fmt.Fprintln(b.Socket, string(data))
}

func (b *Broker) resizeCbArray() {
	var i int
	cut := false
	for i = len(b.Callbacks); i > 0; i-- {
		if b.Callbacks[i-1] != nil {
			break
		}
		cut = true
	}
	if cut {
		b.Callbacks = b.Callbacks[0:i]
		if b.cbFree > i {
			b.cbFree = i
		}
	}
}
