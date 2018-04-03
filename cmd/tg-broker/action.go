package main

import (
	"net"

	"github.com/hamcha/clessy/tg"
)

func executeClientCommand(action tg.ClientCommand, client net.Conn) {
	switch action.Type {
	case tg.CmdSendTextMessage:
		data := *(action.TextMessageData)
		api.SendTextMessage(data)
	case tg.CmdGetFile:
		data := *(action.FileRequestData)
		api.GetFile(data, client, *action.Callback)
	case tg.CmdSendPhoto:
		data := *(action.PhotoData)
		api.SendPhoto(data)
	case tg.CmdForwardMessage:
		data := *(action.ForwardMessageData)
		api.ForwardMessage(data)
	case tg.CmdSendChatAction:
		data := *(action.ChatActionData)
		api.SendChatAction(data)
	}
}
