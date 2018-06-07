package processors

import (
	"Relay-chat/chatserver/model"
	"net"
	"Relay-chat/chatserver/applog"
	"encoding/json"
	"io"
)

type MessageProcessor struct {

}


func (l * MessageProcessor) ProcessMessage(m model.CommonMessage, c net.Conn) {
	_, ok := model.ReadKey(m.Msg.From)
	if !ok {
		applog.Warning.Println("User",m.Msg.From,"not connected")
		WriteMessage(c, m, fail, "chat-server", m.Msg.From)
		return
	}
	v, ok := model.ReadKey(m.Msg.To)
	if !ok {
		applog.Warning.Println("User",m.Msg.To,"not connected")
		WriteMessage(c, m, fail, "chat-server", m.Msg.From)
		return
	}
	encoder := json.NewEncoder(v.C)
	e := encoder.Encode(&m)
	if e != nil && e == io.EOF {
		applog.Warning.Println("Encoder cannot send data")
		model.DeleteFromMap(m.Msg.To)
	}
	applog.Info.Println("From:",m.Msg.From,"To:",m.Msg.To,"Message:",m.Msg.Data)
}

