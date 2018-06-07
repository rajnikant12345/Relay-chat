package processors

import (
	"Relay-chat/chatserver/applog"
	"Relay-chat/chatserver/model"
	"encoding/json"
	"io"
	"net"
)

type KeyExchangeProcessor struct {
}

func (l *KeyExchangeProcessor) ProcessMessage(m model.CommonMessage, c net.Conn) {
	applog.Info.Println("Processig Key Exchange Message")
	_, ok := model.ReadKey(m.KeyExchg.From)
	if !ok {
		applog.Warning.Println("User", m.KeyExchg.From, "not connected")
		WriteMessage(c, m, fail, "chat-server", m.KeyExchg.From)
		return
	}
	val, ok := model.ReadKey(m.KeyExchg.To)
	if !ok {
		applog.Warning.Println("User", m.KeyExchg.To, "not connected")
		WriteMessage(c, m, fail, "chat-server", m.KeyExchg.From)
		return
	}
	encoder := json.NewEncoder(val.Con)
	e := encoder.Encode(&m)
	if e != nil && e == io.EOF {
		applog.Warning.Println("Encoder cannot send data")
		model.DeleteFromMap(m.KeyExchg.To)
	}
	applog.Info.Println("Key Exchange", "From:", m.KeyExchg.From, "To:", m.KeyExchg.To, "Message:", m.KeyExchg.Key)
}
