package processors

import (
	"Relay-chat/chatserver/model"
	"encoding/json"
	"io"
	"net"
	"Relay-chat/chatserver/applog"
)

const fail = "-1"
const loginsuccess = "0"


func GetProcessor(message model.CommonMessage) Processor {
	applog.Info.Println("Getting message processor")
	if message.Lgin != nil {
		return &LoginProcessor{}
	}
	if message.Msg != nil {
		return &MessageProcessor{}
	}
	if message.KeyExchg != nil {
		return &KeyExchangeProcessor{}
	}
	applog.Error.Println("Invalid message processor")
	return nil
}

func WriteMessage(c net.Conn, m model.CommonMessage, msg string, from, to string) {
	applog.Info.Println("WriteMessage","From:",from,"To:",to,"Message:",msg)
	encoder := json.NewEncoder(c)
	message := model.Message{}
	message.Data = msg
	message.From = from
	message.To = to
	m.Msg = &message
	e := encoder.Encode(&m)
	if e != nil && e == io.EOF {
		applog.Warning.Println("Encoder cannot send data")
		model.DeleteFromMap(to)
	}
}